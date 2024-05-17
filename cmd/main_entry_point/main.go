package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

const readHeaderTimeout = 5

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())

	httpSrv := &http.Server{
		Addr:              `:8080`,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	cfg := loadDBEnv(ctx, locateDefaultEnvFile())

	serverWorkers := &sync.WaitGroup{}

	_ = mustInitDBConnection(ctx, cfg, serverWorkers)
	_ = initMemcachedConn(ctx, cfg, serverWorkers)
	upPrometheusExporter(ctx, serverWorkers)

	serverMetrics := registerPrometheusMetrics()

	http.HandleFunc("/", createWelcomeHandler(serverMetrics))
	log.Printf("Starting qr-code-extractor server on addr: %s", httpSrv.Addr)

	go func() {
		if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	gracefulShutdown(ctx, cancelFunc, httpSrv, serverWorkers)

	log.Printf(`qr-code-extractor finished job, bye`)
}

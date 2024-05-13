package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	const readHeaderTimeout = 5

	ctx, cancelFunc := context.WithCancel(context.Background())

	httpSrv := &http.Server{
		Addr:              `:8080`,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	cfg := loadDBEnv(ctx, locateDefaultEnvFile())

	serverWorkers := &sync.WaitGroup{}

	_ = mustInitDBConnection(ctx, cfg, serverWorkers)
	_ = initMemcachedConn(ctx, cfg, serverWorkers)

	http.HandleFunc("/", welcomeHandler)
	log.Printf("Starting qr-code-extractor server on addr: %s", httpSrv.Addr)

	go func() {
		if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	gracefulShutdown(ctx, cancelFunc, httpSrv, serverWorkers)

	log.Printf(`qr-code-extractor finished job, bye`)
}

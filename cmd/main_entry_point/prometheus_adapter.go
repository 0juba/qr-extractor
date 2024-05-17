package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func upPrometheusExporter(ctx context.Context, serverWorkersWG *sync.WaitGroup) {
	const exporterAddr = `:62112`

	router := http.NewServeMux()
	handler := promhttp.Handler()

	router.HandleFunc(`/metrics`, handler.ServeHTTP)

	server := &http.Server{
		Addr:              exporterAddr,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	serverWorkersWG.Add(1)

	go func() {
		log.Printf(`up prometheus metrics on %s`, exporterAddr)

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	go func() {
		defer serverWorkersWG.Done()
		<-ctx.Done()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Print(err)
		} else {
			log.Print(`gracefully stopped prometheus`)
		}
	}()
}

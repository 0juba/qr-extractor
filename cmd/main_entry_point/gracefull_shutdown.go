package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func gracefulShutdown(ctx context.Context, cancelFunc context.CancelFunc, httpSrv *http.Server,
	serverWorkers *sync.WaitGroup,
) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
	<-sigint

	cancelFunc()

	const gracefulShutdownReleaseTimeout = 10

	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, time.Second*gracefulShutdownReleaseTimeout)
	defer shutdownRelease()

	if err := httpSrv.Shutdown(shutdownCtx); err != nil {
		log.Panicf(`http shutdown err %s`, err)
	}

	waitor := make(chan bool)
	go func() {
		serverWorkers.Wait()
		waitor <- true
	}()

	select {
	case <-shutdownCtx.Done():
	case <-waitor:
	}
}

package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5"
)

func withHTTPSrv(srv *http.Server) func(ctx context.Context) {
	return func(ctx context.Context) {
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
	}
}

func withDBConn(dbConn *pgx.Conn) func(ctx context.Context) {
	return func(ctx context.Context) {
		err := dbConn.Close(ctx)
		if err != nil {
			log.Println(err)
		}
	}
}

func gracefulShutdown(ctx context.Context, opts ...func(ctx context.Context)) chan struct{} {
	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		for _, o := range opts {
			o(ctx)
		}

		close(idleConnectionsClosed)
	}()

	return idleConnectionsClosed
}

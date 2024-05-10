package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

func main() {
	const readHeaderTimeout = 5

	ctx := context.Background()

	httpSrv := http.Server{
		Addr:              `:8080`,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	dbConn := mustInitDBConnection(ctx)

	idleConnectionsClosed := gracefulShutdown(ctx, withDBConn(dbConn), withHTTPSrv(&httpSrv))

	http.HandleFunc("/", welcomeHandler)
	log.Printf("Starting qr-code-extractor server on addr: %s", httpSrv.Addr)

	if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-idleConnectionsClosed

	log.Printf(`qr-code-extractor finished job, bye`)
}

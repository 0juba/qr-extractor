package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func welcomeHandler(resp http.ResponseWriter, req *http.Request) {
	log.Printf("HTTP %s %s%s from %s\n", req.Method, req.Host, req.URL, req.RemoteAddr)

	if req.URL.Path != "/" {
		http.Error(resp, "Not Found", http.StatusNotFound)

		return
	}

	resp.Header().Set(`Content-Type`, `application/json`)

	err := json.NewEncoder(resp).Encode(map[string]string{
		`app`: `qr-code-extractor`,
		`v`:   `1.0`,
	})
	if err != nil {
		log.Printf(`HTTP error occurred for remote addr %s\n`, req.RemoteAddr)

		resp.WriteHeader(http.StatusInternalServerError)
		_, _ = resp.Write([]byte(`{"error": "internal server error"}`))
	} else {
		resp.WriteHeader(http.StatusOK)
	}
}

func main() {
	const readHeaderTimeout = 5

	httpSrv := http.Server{
		Addr:              `:8080`,
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := httpSrv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}

		close(idleConnectionsClosed)
	}()

	http.HandleFunc("/", welcomeHandler)
	log.Printf("Starting qr-code-extractor server on addr: %s", httpSrv.Addr)

	if err := httpSrv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-idleConnectionsClosed

	log.Printf(`qr-code-extractor finished job, bye`)
}

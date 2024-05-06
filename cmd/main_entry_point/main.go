package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
)

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("HTTP %s %s%s from %s\n", r.Method, r.Host, r.URL, r.RemoteAddr)

	if r.URL.Path != "/" {
		http.Error(w, "Not Found", http.StatusNotFound)

		return
	}

	w.Header().Set(`Content-Type`, `application/json`)

	err := json.NewEncoder(w).Encode(map[string]string{
		`app`: `qr-code-extractor`,
		`v`:   `1.0`,
	})
	if err != nil {
		log.Printf(`HTTP error occurred for remote addr %s\n`, r.RemoteAddr)

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "internal server error"}`))
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	httpSrv := http.Server{
		Addr: `:8080`,
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
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-idleConnectionsClosed

	log.Printf(`qr-code-extractor finished job, bye bye`)
}

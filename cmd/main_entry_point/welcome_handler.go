package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func welcomeHandler(resp http.ResponseWriter, req *http.Request) {
	log.Printf("HTTP %s %s%s from %s\n", req.Method, req.Host, req.URL, req.RemoteAddr)

	if req.URL.Path != "/" {
		http.Error(resp, "Not Found", http.StatusNotFound)

		return
	}

	resp.Header().Set(`Content-Type`, `application/json`)

	rawBody, err := json.Marshal(map[string]string{
		`app`: `qr-code-extractor`,
		`v`:   `1.0`,
	})
	if err != nil {
		log.Printf(`HTTP error occurred for remote addr %s\n`, req.RemoteAddr)

		resp.WriteHeader(http.StatusInternalServerError)
		_, _ = resp.Write([]byte(`{"error": "internal server error"}`))
	} else {
		resp.WriteHeader(http.StatusOK)

		_, err := resp.Write(rawBody)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			_, _ = resp.Write([]byte(`{"error": "internal server error"}`))
		}
	}
}

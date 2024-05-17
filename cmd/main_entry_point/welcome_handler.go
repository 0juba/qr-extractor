package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func createWelcomeHandler(serverMetrics *metrics) func(resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		var finalHTTPCode int

		if serverMetrics != nil {
			timer := prometheus.NewTimer(*serverMetrics.ts)
			defer func() {
				serverMetrics.rps.With(prometheus.Labels{"code": strconv.Itoa(finalHTTPCode)}).Inc()
				timer.ObserveDuration()
			}()
		}

		log.Printf("HTTP %s %s%s from %s\n", req.Method, req.Host, req.URL, req.RemoteAddr)

		if req.URL.Path != "/" {
			http.Error(resp, "Not Found", http.StatusNotFound)
			finalHTTPCode = http.StatusNotFound

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
			finalHTTPCode = http.StatusInternalServerError

			_, _ = resp.Write([]byte(`{"error": "internal server error"}`))
		} else {
			resp.WriteHeader(http.StatusOK)
			finalHTTPCode = http.StatusOK

			_, err := resp.Write(rawBody)
			if err != nil {
				resp.WriteHeader(http.StatusInternalServerError)
				finalHTTPCode = http.StatusInternalServerError

				_, _ = resp.Write([]byte(`{"error": "internal server error"}`))
			}
		}
	}
}

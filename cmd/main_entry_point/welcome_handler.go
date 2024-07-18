package main

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
)

func createWelcomeHandler(
	serverMetrics *metrics,
	memcacheClient *memcache.Client,
	pgConn *pgx.Conn,
) func(
	resp http.ResponseWriter, req *http.Request) {
	return func(resp http.ResponseWriter, req *http.Request) {
		var finalHTTPCode int

		ctx := req.Context()

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

		waitor := &sync.WaitGroup{}
		handleQueryParams(ctx, waitor, req, memcacheClient, pgConn)

		waitor.Add(1)

		go func() {
			defer waitor.Done()

			finalHTTPCode = handleMainReq(resp, req)
		}()

		waitor.Wait()
	}
}

func handleQueryParams(ctx context.Context, waitor *sync.WaitGroup, req *http.Request, memcacheClient *memcache.Client,
	pgConn *pgx.Conn,
) {
	switch {
	case req.URL.Query().Has("cpu-bound"):
		waitor.Add(1)

		go func() {
			defer waitor.Done()

			v, err := strconv.ParseInt(req.URL.Query().Get("cpu-bound"), 10, 64)
			if err == nil {
				onCPUBound(int(v))
			}
		}()
	case req.URL.Query().Has("cache-w"):
		waitor.Add(1)

		go func() {
			defer waitor.Done()
			onCacheWrite(memcacheClient)
		}()
	case req.URL.Query().Has("cache-r"):
		waitor.Add(1)

		go func() {
			defer waitor.Done()
			onCacheRead(memcacheClient)
		}()
	case req.URL.Query().Has("db-r"):
		waitor.Add(1)

		go func() {
			defer waitor.Done()
			onDPCPUBound(ctx, pgConn)
		}()
	}
}

func handleMainReq(resp http.ResponseWriter, req *http.Request) int {
	var finalHTTPCode int

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

	return finalHTTPCode
}

func onCPUBound(count int) {
	v := rand.ExpFloat64()
	for range count {
		sha512.Sum512_256([]byte(fmt.Sprintf("%f", v)))
	}

	log.Printf("cpu bound req")
}

func onCacheWrite(memcachedClient *memcache.Client) {
	k, v := rand.ExpFloat64(), rand.ExpFloat64()
	rawV := fmt.Sprintf("%f", v)

	err := memcachedClient.Set(&memcache.Item{
		Key:   fmt.Sprintf("random-key-%f", k),
		Value: []byte(rawV),
	})
	if err != nil {
		log.Printf("cannot write cache: %s", err)
	} else {
		log.Printf("write cache successfully")
	}
}

func onCacheRead(memcachedClient *memcache.Client) {
	v := rand.ExpFloat64()
	item, err := memcachedClient.Get(fmt.Sprintf("random-key-%f", v))

	if err != nil {
		log.Printf("cannot read from cache: %s", err)
	} else {
		log.Printf("read from cache: %v", item)
	}
}

func onDPCPUBound(ctx context.Context, pgConn *pgx.Conn) {
	resultSet, err := pgConn.Query(ctx,
		"SELECT md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea), "+
			"md5(random()::text), sha512(random()::text::bytea)")
	if err != nil {
		log.Printf("read from db err: %s", err)
	}

	defer resultSet.Close()

	if resultSet.Err() != nil {
		log.Printf("db internal err: %s", resultSet.Err())
	} else {
		log.Printf("got query result")
	}
}

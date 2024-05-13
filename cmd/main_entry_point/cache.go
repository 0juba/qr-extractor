package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

func initMemcachedConn(ctx context.Context, cfg *Env, serverWorkers *sync.WaitGroup) *memcache.Client {
	dsn := net.JoinHostPort(cfg.Memcached.Host, cfg.Memcached.Port)

	memcachedClient := memcache.New(dsn)
	log.Printf(`up memcached conn on %s`, dsn)

	go memcachedPing(ctx, memcachedClient, serverWorkers)
	go memcachedOnCtxStop(ctx, memcachedClient)

	return memcachedClient
}

func memcachedOnCtxStop(ctx context.Context, mc *memcache.Client) {
	<-ctx.Done()

	err := mc.Close()
	if err != nil {
		log.Printf(`cannot gracefully close memcached connection due to %s`, err)
	} else {
		log.Printf(`memcached conn gracefully stopped`)
	}
}

func memcachedPing(ctx context.Context, memcachedClient *memcache.Client, wg *sync.WaitGroup) {
	const cachePingIntervalSec = 5

	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(cachePingIntervalSec * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf(`gracefully stopped memcached ping worker`)

			return
		case <-ticker.C:
			err := memcachedClient.Ping()
			if err != nil {
				log.Printf(`cannot ping memcached due to %s`, err)
			} else {
				log.Printf(`successfully ping memcached`)
			}
		}
	}
}

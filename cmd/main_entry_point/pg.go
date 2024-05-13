package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
)

func mustInitDBConnection(ctx context.Context, cfg *Env, serverWorkers *sync.WaitGroup) *pgx.Conn {
	conn, err := pgx.Connect(ctx, mustBuildDSN(cfg))
	if err != nil {
		panic(err)
	}

	go dbWarmer(ctx, conn, serverWorkers)
	go dbOnCtxDone(ctx, conn)

	log.Printf(`up pg connection on %s:%s`, cfg.DB.Host, cfg.DB.Port)

	return conn
}

func dbOnCtxDone(ctx context.Context, conn *pgx.Conn) {
	<-ctx.Done()

	err := conn.Close(ctx)
	if err != nil {
		log.Printf(`cannot gracefully close pg conn due to %s`, err)
	} else {
		log.Printf(`pg conn gracefully stopped`)
	}
}

func dbWarmer(ctx context.Context, pgConn *pgx.Conn, wg *sync.WaitGroup) {
	const pgPingInterval = 5

	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(pgPingInterval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf(`gracefully stopped db ping worker`)

			return
		case <-ticker.C:
			err := pgConn.Ping(ctx)
			if err != nil {
				log.Printf(`pg ping failed due to %v`, err)

				return
			}

			log.Printf(`pg db ping successfully done`)
		}
	}
}

func mustBuildDSN(cfg *Env) string {
	return fmt.Sprintf(
		`postgres://%s:%s@%s/%s`,
		cfg.DB.User,
		cfg.DB.Pwd,
		net.JoinHostPort(cfg.DB.Host, cfg.DB.Port),
		cfg.DB.Name,
	)
}

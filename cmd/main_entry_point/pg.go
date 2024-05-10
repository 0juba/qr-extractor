package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	envcfg "github.com/sethvargo/go-envconfig"
)

type Env struct {
	DBHost string `env:"DB_HOST"`
	DBPort string `env:"DB_PORT"`
	DBUser string `env:"DB_USER"`
	DBPwd  string `env:"DB_PASSWORD"`
	DBName string `env:"DB_NAME"`
}

func mustInitDBConnection(ctx context.Context) *pgx.Conn {
	cfg := loadDBEnv(ctx, locateDefaultEnvFile())

	conn, err := pgx.Connect(ctx, mustBuildDSN(cfg))
	if err != nil {
		panic(err)
	}

	go dbWarmer(ctx, conn)

	log.Printf(`up pg connection on %s:%s`, cfg.DBHost, cfg.DBPort)

	return conn
}

func dbWarmer(ctx context.Context, c *pgx.Conn) {
	const pgPingInterval = 5

	t := time.NewTicker(pgPingInterval * time.Second)
	for range t.C {
		err := c.Ping(ctx)
		if err != nil {
			log.Printf(`pg ping failed due to %v`, err)

			return
		}

		log.Printf(`pg db ping successfully done`)
	}
}

func mustBuildDSN(cfg Env) string {
	return fmt.Sprintf(
		`postgres://%s:%s@%s/%s`, cfg.DBUser, cfg.DBPwd, net.JoinHostPort(cfg.DBHost, cfg.DBPort), cfg.DBName)
}

func locateDefaultEnvFile() string {
	relative := `./.env`

	envPath, err := filepath.Abs(relative)
	if err != nil {
		panic(err)
	}

	return envPath
}

func loadDBEnv(ctx context.Context, envFile string) Env {
	envFileFD, err := os.Open(envFile)
	if err != nil {
		log.Fatal(`cannot open envfile due to`, err)
	}

	defer func(fd *os.File) {
		_ = fd.Close()
	}(envFileFD)

	c := Env{}
	_ = godotenv.Load(envFile)
	_ = envcfg.Process(ctx, &c)

	return c
}

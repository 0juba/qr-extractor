package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	envcfg "github.com/sethvargo/go-envconfig"
)

type (
	Env struct {
		DB        *DB        `env:", prefix=DB_"`
		Memcached *Memcached `env:", prefix=CACHE_"`
	}

	Memcached struct {
		Host string `env:"HOST"`
		Port string `env:"PORT"`
	}

	DB struct {
		Host string `env:"HOST"`
		Port string `env:"PORT"`
		User string `env:"USER"`
		Pwd  string `env:"PASSWORD"`
		Name string `env:"NAME"`
	}
)

func locateDefaultEnvFile() string {
	relative := `./.env`

	envPath, err := filepath.Abs(relative)
	if err != nil {
		panic(err)
	}

	return envPath
}

func loadDBEnv(ctx context.Context, envFile string) *Env {
	envFileFD, err := os.Open(envFile)
	if err != nil {
		log.Fatal(`cannot open envfile due to`, err)
	}

	defer func(fd *os.File) {
		_ = fd.Close()
	}(envFileFD)

	c := &Env{}
	_ = godotenv.Load(envFile)
	_ = envcfg.Process(ctx, c)

	return c
}

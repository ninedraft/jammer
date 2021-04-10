package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ninedraft/gemax/gemax"
	"github.com/ninedraft/jammer/config"
	"github.com/ninedraft/jammer/counter"
	"github.com/ninedraft/jammer/middleware"
)

func main() {
	var cfg = &config.Config{}
	cfg.Flag(flag.CommandLine)
	flag.Parse()

	var tlsConfig, errTLS = buildTLS(cfg)
	if errTLS != nil {
		log.Printf("%v", errTLS)
		os.Exit(1)
	}

	var blog = middleware.Use(
		(&gemax.FileSystem{
			Prefix: cfg.ContentDir,
			FS:     os.DirFS(cfg.ContentDir),
			Logf:   log.Printf,
		}).Serve,
		counter.New("/stats").Middleware,
	)

	var server = gemax.Server{
		Addr:    cfg.Addr,
		Hosts:   cfg.Hosts,
		Logf:    log.Printf,
		Handler: blog,
	}

	var ctx, cancel = signal.NotifyContext(context.Background())
	defer cancel()

	var errServe = server.ListenAndServe(ctx, tlsConfig)
	if errServe != nil {
		log.Printf("serving: %v", errServe)
	}
}

func buildTLS(cfg *config.Config) (*tls.Config, error) {
	var cert, errLoad = tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if errLoad != nil {
		return nil, fmt.Errorf("loading TLS cert: %w", errLoad)
	}
	return &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
	}, nil
}

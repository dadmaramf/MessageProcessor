package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"messageprocessor/internal/app"
	"messageprocessor/internal/config"
	messagesender "messageprocessor/internal/services/message_sender"
	storagep "messageprocessor/internal/storage/postgres"
	"messageprocessor/pkg/kafka"
	"messageprocessor/pkg/postgres"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	if err := RunApp(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func RunApp(ctx context.Context, write io.Writer) error {

	log := slog.New(slog.NewJSONHandler(write, &slog.HandlerOptions{Level: slog.LevelDebug}))

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		return err
	}

	storage := storagep.NewMessageStorage(db)

	producer, err := kafka.NewSyncProducer(cfg)

	if err != nil {
		return err
	}
	service := messagesender.New(storage, producer, log)

	server := app.NewServer(log, cfg, service)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(cfg.Http.Host, cfg.Http.Port),
		Handler: server,
	}

	go func() {
		log.Info("Listening on", "Address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	service.StartProcessingMessage(ctx, time.Second*10)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

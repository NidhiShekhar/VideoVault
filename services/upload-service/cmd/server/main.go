package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	apphttp "github.com/nidhi/video-vault/services/upload-service/internal/http"
	"github.com/nidhi/video-vault/services/upload-service/internal/config"
	"github.com/nidhi/video-vault/services/upload-service/internal/repository/memory"
	"github.com/nidhi/video-vault/services/upload-service/internal/service"
)

func main() {
	configPath := flag.String("config", "config/config.yaml", "Path to service config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	logger := newLogger(cfg.Log.Level)
	slog.SetDefault(logger)

	repo := memory.NewVideoRepository()
	videoService := service.NewVideoService(repo)
	uploadHandler := apphttp.NewUploadHandler(videoService)
	healthHandler := apphttp.NewHealthHandler()
	router := apphttp.NewRouter(uploadHandler, healthHandler)

	srv := &http.Server{
		Addr:              ":" + cfg.Server.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		slog.Info("starting upload-service",
			"port", cfg.Server.Port,
			"env", cfg.Service.Env,
		)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server failed", "error", err)
			stop()
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		return
	}

	slog.Info("server stopped gracefully")
}

func newLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch level {
	case "debug":
		lvl = slog.LevelDebug
	case "warn":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
}

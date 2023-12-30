package main

import (
    "log/slog"
    "os"
    "url_shortener/internal/config"
)

func main() {
    cfg := config.MustLoad()
    logger := setupLogger(cfg.Env)
    logger.Info("Starting url-shortener")
    logger.Debug("Debug level enabled")
}

func setupLogger(env string) *slog.Logger {
    var log *slog.Logger
    switch env {
    case "local":
        log = slog.New(
            slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case "dev":
        log = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
        )
    case "prod":
        log = slog.New(
            slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
        )
    }
    return log
}

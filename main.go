package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"eolscream/pkg"
	"github.com/caarlos0/env/v10"
)

type Config struct {
	SlackWebhookURL string `env:"SLACK_WEBHOOK_URL" envDefault:""` // Assuming no default value for Slack webhook URL
	CatalogueFile   string `env:"CATALOGUE_FILE" envDefault:"catalogue.json"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"INFO"`
}

func main() {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		panic(fmt.Errorf("unable to parse config: %w", err))
	}

	var notifier pkg.Notifier
	switch config.SlackWebhookURL {
	case "":
		notifier = pkg.NewNilNotifier()
	default:
		notifier = pkg.NewSlackNotifier(config.SlackWebhookURL)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: calcLevel(config.LogLevel),
	}))

	checker := pkg.NewCatalogueChecker(pkg.CatalogueCheckerOptions{
		Path:     config.CatalogueFile,
		Client:   pkg.NewEndOfLifeDateClient(http.DefaultClient),
		Notifier: notifier,
		Logger:   logger,
	})

	err := checker.NotifyNearEndOfLife()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

func calcLevel(levelAsString string) slog.Level {
	var level slog.Level

	switch levelAsString {
	case "Debug":
		level = slog.LevelDebug
	case "Warn":
		level = slog.LevelWarn
	case "Error":
		level = slog.LevelError
	case "Info":
	default:
		level = slog.LevelInfo
	}

	return level
}

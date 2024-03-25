package main

import (
	"log/slog"
	"net/http"
	"os"

	"eolscream/pkg"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	checker := pkg.NewCatalogueChecker(pkg.CatalogueCheckerOptions{
		Path:     "catalogue.json",
		Client:   pkg.NewEndOfLifeDateClient(http.DefaultClient),
		Notifier: pkg.NewNilNotifier(),
		Logger:   logger,
	})

	err := checker.NotifyNearEndOfLife()
	if err != nil {
		logger.Error(err.Error())
		return
	}
}

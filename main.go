package main

import (
	"log/slog"
	"net/http"
	"os"

	"eolscream/pkg"
)

func main() {
	checker := pkg.NewCatalogueChecker(pkg.CatalogueCheckerOptions{
		Path:     "catalogue.json",
		Client:   pkg.NewEndOfLifeDateClient(http.DefaultClient),
		Notifier: pkg.NewNilNotifier(),
		Logger:   slog.New(slog.NewJSONHandler(os.Stdout, nil)),
	})

	err := checker.NotifyNearEndOfLife()
	if err != nil {
		return
	}
}

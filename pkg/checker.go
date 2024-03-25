package pkg

import (
	"fmt"
	"log/slog"
)

type CatalogueChecker struct {
	opts CatalogueCheckerOptions
}

type CatalogueCheckerOptions struct {
	Path     string
	Client   EndOfLifeClient
	Notifier Notifier
	Logger   *slog.Logger
}

func NewCatalogueChecker(opts CatalogueCheckerOptions) *CatalogueChecker {
	return &CatalogueChecker{opts: opts}
}

func (c *CatalogueChecker) NotifyNearEndOfLife() error {
	logger := c.opts.Logger

	logger.Debug("attempting to load catalogue from path", slog.String("path", c.opts.Path))
	catalogue, err := NewCatalogueFromFile(c.opts.Path)
	if err != nil {
		return fmt.Errorf("CatalogueChecker@NotifyNearEndOfLife: unable to load catalogue: %w", err)
	}

	for _, product := range *catalogue {
		logger.Debug("attempting to fetch release info for product", slog.String("product", product.Name), slog.String("version", product.Version))
		releaseInfo, err := c.opts.Client.FetchReleaseInfo(product)
		if err != nil {
			logger.Error("CatalogueChecker@NotifyNearEndOfLife: unable to fetch release info for product: %w", "error", err)
			return nil
		}

		nearEol, err := c.inspectReleaseInfo(logger, product, releaseInfo, err)
		if err != nil {
			return err
		}

		if nearEol {
			c.opts.Logger.Debug(
				"product is near eol, sending notification",
				slog.String("product", product.Name),
				slog.String("version", product.Version),
				slog.String("eol", releaseInfo.EndOfLifeDate),
				slog.String("release", releaseInfo.ReleaseDate),
			)

			c.opts.Notifier.Notify(product, *releaseInfo)
		}
	}

	return nil
}

func (c *CatalogueChecker) inspectReleaseInfo(logger *slog.Logger, product Product, releaseInfo *ReleaseInfo, err error) (bool, error) {
	logger.Debug(
		"checking eol within release info",
		slog.String("product", product.Name),
		slog.String("eol", releaseInfo.EndOfLifeDate),
		slog.String("release", releaseInfo.ReleaseDate),
	)

	nearEol, err := releaseInfo.IsEndOfLifeDateNear()
	if err != nil {
		return false, fmt.Errorf("CatalogueChecker@NotifyNearEndOfLife: unable to check is near eol: %w", err)
	}

	return nearEol, nil
}

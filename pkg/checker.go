package pkg

import "log/slog"

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
	catalogue, err := NewCatalogueFromFile(c.opts.Path)
	if err != nil {
		return err
	}

	for _, product := range *catalogue {
		releaseInfo, err := c.opts.Client.FetchReleaseInfo(product)
		if err != nil {
			return err
		}

		nearEol, err := releaseInfo.IsEndOfLifeDateNear()
		if err != nil {
			return err
		}

		if nearEol {
			c.opts.Logger.Info(
				"Product is at EOL",
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

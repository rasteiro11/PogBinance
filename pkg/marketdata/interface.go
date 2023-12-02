package marketdata

import (
	"context"

	"github.com/rasteiro11/PogBinance/models"
)

type MarketDataProvider interface {
	Start(ctx context.Context, currency ...string) (<-chan models.MarketData, error)
	Close(ctx context.Context) error
}

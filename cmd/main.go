package main

import (
	"context"
	binanceProvider "github.com/rasteiro11/PogBinance/pkg/marketdata/binance"

	"github.com/rasteiro11/PogCore/pkg/logger"
)

func main() {
	ctx := context.Background()

	binanceWs := binanceProvider.NewMarketDataProvider()

	c, err := binanceWs.Start(ctx, "maticusdt", "dotusdt", "solusdt", "avaxusdt")
	if err != nil {
		logger.Of(ctx).Fatalf("[main] binanceWs.Start() returned error: %+v\n", err)
	}

	for md := range c {
		logger.Of(ctx).Debugf("MARKET DATA: %+v\n", md)
	}

	for {
	}
}

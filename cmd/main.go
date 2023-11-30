package main

import (
	"context"
	"log"
	binanceProvider "pogbinance/pkg/marketdata/binance"

	"github.com/rasteiro11/PogCore/pkg/logger"
)

func main() {
	ctx := context.Background()

	binanceWs := binanceProvider.NewMarketDataProvider()

	c, err := binanceWs.Start(ctx, "usdtbrl")
	if err != nil {
		logger.Of(ctx).Fatalf("[main] binanceWs.Start() returned error: %+v\n", err)
	}

	for md := range c {
		log.Printf("RECEIVED MARKET DATA: %+v\n", md)
	}

	for {
	}
}

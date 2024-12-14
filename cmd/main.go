package main

import (
	"log"

	"hft-trading-system/internal/services"
)

func main() {
	const (
		binanceWSURL = "wss://stream.binance.com:9443/ws"
		krakenWSURL  = "wss://ws.kraken.com"
	)

	// Connect Binance
	binanceClient, err := services.NewBinanceClient(binanceWSURL)
	if err != nil {
		log.Fatalf("Failed to connect to Binance: %v", err)
	}

	// Connect Kraken
	krakenClient, err := services.NewKrakenClient(krakenWSURL)
	if err != nil {
		log.Fatalf("Failed to connect to Kraken: %v", err)
	}

	// Creating a service with both clients
	marketDataService := services.NewMarketDataService(binanceClient, krakenClient)

	// Subscribe to the tickers and start processing
	symbol := "btcusdt"
	marketDataService.Start(symbol)

	select {}
}

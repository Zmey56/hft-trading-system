package main

import (
	"log"

	"hft-trading-system/internal/services"
)

func main() {
	const (
		binanceWSURL = "wss://stream.binance.com:9443/ws"
		krakenWSURL  = "wss://ws.kraken.com"
		threshold    = 10.0 // Setting the minimum difference for arbitration
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
	// marketDataService := services.NewMarketDataService(binanceClient, krakenClient)

	// Creating and launching an arbitration service
	arbitrageService := services.NewArbitrageService(binanceClient, krakenClient, threshold)
	symbol := "btcusdt"
	arbitrageService.Start(symbol)

	// Subscribe to the tickers and start processing
	// symbol := "btcusdt"
	// marketDataService.Start(symbol)

	select {}
}

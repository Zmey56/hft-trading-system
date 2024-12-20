package main

import (
	"log"

	"hft-trading-system/internal/services"
)

func main() {
	const krakenWSURL = "wss://ws.kraken.com"
	const binanceWSURL = "wss://stream.binance.com:9443/ws"

	// Подключаем Kraken
	krakenClient, err := services.NewKrakenClient(krakenWSURL)
	if err != nil {
		log.Fatalf("Failed to connect to Kraken: %v", err)
	}

	// Подключаем Binance
	binanceClient, err := services.NewBinanceClient(binanceWSURL)
	if err != nil {
		log.Fatalf("Failed to connect to Binance: %v", err)
	}

	// Подписываемся на тикеры
	symbol := "btcusdt"

	if err := krakenClient.SubscribeToTicker(symbol); err != nil {
		log.Fatalf("Failed to subscribe to Kraken ticker: %v", err)
	}
	if err := binanceClient.SubscribeToTicker(symbol); err != nil {
		log.Fatalf("Failed to subscribe to Binance ticker: %v", err)
	}

	// Создаём и запускаем сервис арбитража
	arbitrageService := services.NewArbitrageService(binanceClient, krakenClient)
	arbitrageService.Start(symbol)
}

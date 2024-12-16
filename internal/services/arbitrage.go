package services

import (
	"encoding/json"
	"log"
)

// BinanceTickerMessage Defining structures for messages processing
type BinanceTickerMessage struct {
	Stream string `json:"stream"`
	Data   struct {
		BidPrice string `json:"b"`
		AskPrice string `json:"a"`
	} `json:"data"`
}

type KrakenTickerMessage struct {
	ChannelID int      `json:"channelID"`
	Pair      string   `json:"pair"`
	Bid       []string `json:"b"`
	Ask       []string `json:"a"`
}

// ArbitrageService is responsible for the execution of the arbitration strategy.
type ArbitrageService struct {
	BinanceClient ExchangeClient
	KrakenClient  ExchangeClient
	Threshold     float64
}

// NewArbitrageService creates a new Arbitrage service
func NewArbitrageService(binanceClient, krakenClient ExchangeClient, threshold float64) *ArbitrageService {
	return &ArbitrageService{
		BinanceClient: binanceClient,
		KrakenClient:  krakenClient,
		Threshold:     threshold,
	}
}

// Start starts the arbitrage service
func (a *ArbitrageService) Start(symbol string) {
	binanceChan := make(chan float64)
	krakenChan := make(chan float64)

	go a.readPrices(a.BinanceClient, binanceChan, "Binance")
	go a.readPrices(a.KrakenClient, krakenChan, "Kraken")

	for {
		binancePrice := <-binanceChan
		krakenPrice := <-krakenChan
		a.evaluateArbitrage(binancePrice, krakenPrice)
	}
}

// readPrices обрабатывает сообщения WebSocket и передаёт цены в канал
func (a *ArbitrageService) readPrices(client ExchangeClient, priceChan chan float64, exchange string) {
	for {
		message, err := client.ReadMessage()
		if err != nil {
			log.Printf("[%s] Error reading message: %v", exchange, err)
			break
		}

		switch exchange {
		case "Binance":
			var binanceMsg BinanceTickerMessage
			if err := json.Unmarshal(message, &binanceMsg); err != nil {
				log.Printf("[%s] Failed to parse message: %v", exchange, err)
				continue
			}
			// Вычисляем среднюю цену
			bidPrice := parsePrice(binanceMsg.Data.BidPrice)
			askPrice := parsePrice(binanceMsg.Data.AskPrice)
			price := (bidPrice + askPrice) / 2
			priceChan <- price

		case "Kraken":
			var krakenMsg KrakenTickerMessage
			if err := json.Unmarshal(message, &krakenMsg); err != nil {
				log.Printf("[%s] Failed to parse message: %v", exchange, err)
				continue
			}
			// Вычисляем среднюю цену
			bidPrice := parsePrice(krakenMsg.Bid[0])
			askPrice := parsePrice(krakenMsg.Ask[0])
			price := (bidPrice + askPrice) / 2
			priceChan <- price
		}
	}
}

// evaluateArbitrage проверяет возможности арбитража
func (a *ArbitrageService) evaluateArbitrage(binancePrice, krakenPrice float64) {
	if binancePrice-krakenPrice > a.Threshold {
		log.Printf("Arbitrage opportunity: Buy on Kraken at %.2f, Sell on Binance at %.2f", krakenPrice, binancePrice)
	} else if krakenPrice-binancePrice > a.Threshold {
		log.Printf("Arbitrage opportunity: Buy on Binance at %.2f, Sell on Kraken at %.2f", binancePrice, krakenPrice)
	} else {
		log.Printf("No arbitrage opportunity. Binance: %.2f, Kraken: %.2f", binancePrice, krakenPrice)
	}
}

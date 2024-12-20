package services

import (
	"encoding/json"
	"log"
	"strconv"
)

// BinanceTickerMessage описывает структуру сообщения от Binance
type BinanceTickerMessage struct {
	Stream string `json:"stream"`
	Data   struct {
		BidPrice string `json:"b"`
		AskPrice string `json:"a"`
	} `json:"data"`
}

// KrakenTickerMessage описывает структуру сообщения от Kraken
type KrakenTickerMessage struct {
	Event  string `json:"event"`
	Pair   string `json:"pair"`
	Ticker struct {
		Bid []string `json:"b"`
		Ask []string `json:"a"`
	} `json:"ticker"`
}

// BinanceResponse для обработки сообщений подтверждения подписки
type BinanceResponse struct {
	Result interface{} `json:"result,omitempty"`
	ID     int         `json:"id,omitempty"`
}

// ArbitrageService управляет получением данных с обеих бирж
type ArbitrageService struct {
	BinanceClient ExchangeClient
	KrakenClient  ExchangeClient
}

// NewArbitrageService создаёт новый сервис арбитража
func NewArbitrageService(binanceClient, krakenClient ExchangeClient) *ArbitrageService {
	return &ArbitrageService{
		BinanceClient: binanceClient,
		KrakenClient:  krakenClient,
	}
}

// Start запускает получение данных с обеих бирж
func (a *ArbitrageService) Start(symbol string) {
	binanceChan := make(chan float64)
	krakenChan := make(chan float64)

	// Горутина для Binance
	go func() {
		log.Println("[Binance] Starting price reader")
		a.readPrices(a.BinanceClient, binanceChan, "Binance")
		log.Println("[Binance] Price reader exited")
	}()

	// Горутина для Kraken
	go func() {
		log.Println("[Kraken] Starting price reader")
		a.readPrices(a.KrakenClient, krakenChan, "Kraken")
		log.Println("[Kraken] Price reader exited")
	}()

	// Основной цикл обработки данных
	for {
		select {
		case binancePrice := <-binanceChan:
			log.Printf("[Binance] Average Price: %.2f", binancePrice)
		case krakenPrice := <-krakenChan:
			log.Printf("[Kraken] Average Price: %.2f", krakenPrice)
		}
	}
}

// readPrices читает сообщения с биржи и передаёт средние цены в канал
func (a *ArbitrageService) readPrices(client ExchangeClient, priceChan chan float64, exchange string) {
	defer log.Printf("[%s] Exiting readPrices goroutine", exchange)

	for {
		message, err := client.ReadMessage()
		if err != nil {
			log.Printf("[%s] Error reading message: %v", exchange, err)
			break
		}

		log.Printf("[%s] Raw Message: %s", exchange, message)

		switch exchange {
		case "Binance":
			// Шаг 1: Логирование полученного сообщения
			log.Printf("[Binance] Received Raw Message: %s", message)

			// Шаг 2: Попытка распознать сообщение как подтверждение подписки
			var response BinanceResponse
			if err := json.Unmarshal(message, &response); err == nil && response.ID == 1 {
				log.Println("[Binance] Subscription confirmed")
				continue
			}
			log.Println("[Binance] Message is not a subscription confirmation")

			// Шаг 3: Попытка распарсить сообщение в карту (map)
			var raw map[string]interface{}
			if err := json.Unmarshal(message, &raw); err != nil {
				log.Printf("[Binance] Failed to parse message as map: %v", err)
				continue
			}
			log.Println("[Binance] Message successfully parsed as map")

			// Шаг 4: Проверка наличия ключей "b" (bid) и "a" (ask)
			log.Printf("[Binance] Checking keys 'b' and 'a' in the message")
			bidStr, bidOk := raw["b"].(string)
			askStr, askOk := raw["a"].(string)

			if !bidOk || !askOk {
				log.Printf("[Binance] Missing Bid or Ask in message: %s", message)
				continue
			}
			log.Printf("[Binance] Found keys 'b' and 'a': bid=%s, ask=%s", bidStr, askStr)

			// Шаг 5: Попытка преобразовать строки в числа
			bid := parsePrice(bidStr)
			ask := parsePrice(askStr)
			if bid == 0 || ask == 0 {
				log.Printf("[Binance] Failed to parse Bid or Ask as float: bid=%s, ask=%s", bidStr, askStr)
				continue
			}
			log.Printf("[Binance] Successfully parsed prices: bid=%.2f, ask=%.2f", bid, ask)

			// Шаг 6: Расчёт средней цены
			averagePrice := (bid + ask) / 2
			log.Printf("[Binance] Calculated Average Price: %.2f", averagePrice)

			// Шаг 7: Отправка средней цены в канал
			priceChan <- averagePrice
			log.Println("[Binance] Average Price sent to channel")

		case "Kraken":
			// Проверяем, является ли сообщение системным
			var raw map[string]interface{}
			if err := json.Unmarshal(message, &raw); err == nil {
				if event, ok := raw["event"]; ok {
					if event == "subscriptionStatus" {
						log.Printf("[Kraken] Subscription confirmed: %v", raw)
						continue
					}
					if event != "ticker" {
						log.Printf("[Kraken] Ignoring non-ticker event: %s", event)
						continue
					}
				}
			}

			// Проверяем, является ли сообщение тикером
			if message[0] == '[' {
				var rawArray []interface{}
				if err := json.Unmarshal(message, &rawArray); err != nil {
					log.Printf("[Kraken] Failed to parse ticker message: %v", err)
					continue
				}

				if len(rawArray) < 2 {
					log.Println("[Kraken] Malformed ticker message")
					continue
				}

				tickerData, ok := rawArray[1].(map[string]interface{})
				if !ok {
					log.Println("[Kraken] Ticker data not found")
					continue
				}

				// Парсим bid/ask
				bid, ask := parseKrakenPrices(tickerData)
				if bid > 0 && ask > 0 {
					averagePrice := (bid + ask) / 2
					priceChan <- averagePrice
					log.Printf("[Kraken] Average Price: %.2f", averagePrice)
				}
			}
		}
	}
}

// parseKrakenPrices парсит bid/ask из тикерных данных
func parseKrakenPrices(tickerData map[string]interface{}) (float64, float64) {
	bid := 0.0
	ask := 0.0

	if b, ok := tickerData["b"].([]interface{}); ok && len(b) > 0 {
		bid = parsePrice(b[0].(string))
	}
	if a, ok := tickerData["a"].([]interface{}); ok && len(a) > 0 {
		ask = parsePrice(a[0].(string))
	}

	return bid, ask
}

// parsePrice преобразует строку в float64
func parsePrice(priceStr string) float64 {
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		log.Printf("Failed to parse price: %v", err)
		return 0
	}
	return price
}

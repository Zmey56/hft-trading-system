package services

import (
	"log"

	"github.com/gorilla/websocket"
)

type BinanceClient struct {
	Conn *websocket.Conn
}

// NewBinanceClient creates a new Binance client
func NewBinanceClient(url string) (*BinanceClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &BinanceClient{Conn: conn}, nil
}

// SubscribeToTicker implementation for Binance
func (b *BinanceClient) SubscribeToTicker(symbol string) error {
	message := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{symbol + "@ticker"},
		"id":     1,
	}
	return b.Conn.WriteJSON(message)
}

// ReadMessages implementation for Binance
func (b *BinanceClient) ReadMessages() {
	defer b.Conn.Close()
	for {
		_, message, err := b.Conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from Binance: %v", err)
			break
		}
		log.Printf("Binance message: %s", message)
	}
}

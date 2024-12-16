package services

import (
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

// ReadMessage ReadMessages implementation for Binance
func (b *BinanceClient) ReadMessage() ([]byte, error) {
	_, message, err := b.Conn.ReadMessage()
	return message, err
}

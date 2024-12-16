package services

import (
	"github.com/gorilla/websocket"
)

type KrakenClient struct {
	Conn *websocket.Conn
}

// NewKrakenClient creates a new Kraken client
func NewKrakenClient(url string) (*KrakenClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &KrakenClient{Conn: conn}, nil
}

// SubscribeToTicker implementation for Kraken
func (k *KrakenClient) SubscribeToTicker(symbol string) error {
	if symbol == "btcusdt" {
		symbol = "XBT/USD" // Преобразование символа для Kraken
	}

	message := map[string]interface{}{
		"event": "subscribe",
		"pair":  []string{symbol},
		"subscription": map[string]string{
			"name": "ticker",
		},
	}
	return k.Conn.WriteJSON(message)
}

// ReadMessage ReadMessages implementation for Kraken
func (k *KrakenClient) ReadMessage() ([]byte, error) {
	_, message, err := k.Conn.ReadMessage()
	return message, err
}

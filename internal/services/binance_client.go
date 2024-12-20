package services

import (
	"log"

	"github.com/gorilla/websocket"
)

type BinanceClient struct {
	Conn *websocket.Conn
}

func NewBinanceClient(url string) (*BinanceClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	log.Println("[Binance] Successfully connected to WebSocket")
	return &BinanceClient{Conn: conn}, nil
}

func (b *BinanceClient) SubscribeToTicker(symbol string) error {
	message := map[string]interface{}{
		"method": "SUBSCRIBE",
		"params": []string{symbol + "@ticker"},
		"id":     1,
	}
	log.Printf("[Binance] Sending subscription message: %v", message)
	err := b.Conn.WriteJSON(message)
	if err != nil {
		log.Printf("[Binance] Failed to send subscription message: %v", err)
		return err
	}
	log.Println("[Binance] Subscription message sent successfully")
	return nil
}

func (b *BinanceClient) ReadMessage() ([]byte, error) {
	_, message, err := b.Conn.ReadMessage()
	if err != nil {
		log.Printf("[Binance] Error reading message: %v", err)
		return nil, err
	}
	log.Printf("[Binance] Raw Message: %s", message)
	return message, nil
}

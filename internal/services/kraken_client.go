package services

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type KrakenClient struct {
	Conn *websocket.Conn
}

// NewKrakenClient создает соединение с Kraken WebSocket API
func NewKrakenClient(url string) (*KrakenClient, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	// Настраиваем обработчики ping/pong
	conn.SetPingHandler(func(appData string) error {
		log.Println("[Kraken] Received ping, sending pong...")
		return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
	})

	conn.SetPongHandler(func(appData string) error {
		log.Println("[Kraken] Received pong")
		return nil
	})

	// Регулярно отправляем ping для поддержания соединения
	go func() {
		for {
			time.Sleep(15 * time.Second) // Интервал пинга
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("[Kraken] Failed to send ping: %v", err)
				return
			}
			log.Println("[Kraken] Sent ping")
		}
	}()

	log.Println("[Kraken] Successfully connected to WebSocket")
	return &KrakenClient{Conn: conn}, nil
}

func (k *KrakenClient) SubscribeToTicker(symbol string) error {
	// Используем правильный формат символа
	if symbol == "btcusdt" {
		symbol = "XBT/USD"
	}

	message := map[string]interface{}{
		"event": "subscribe",
		"pair":  []string{symbol},
		"subscription": map[string]string{
			"name": "ticker",
		},
	}

	// Выводим запрос в лог перед отправкой
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("[Kraken] Failed to marshal subscription message: %v", err)
		return err
	}
	log.Printf("[Kraken] Sending subscription message: %s", string(jsonMessage))

	return k.Conn.WriteJSON(message)
}

// ReadMessage читает сообщение из WebSocket
func (k *KrakenClient) ReadMessage() ([]byte, error) {
	_, message, err := k.Conn.ReadMessage()
	return message, err
}

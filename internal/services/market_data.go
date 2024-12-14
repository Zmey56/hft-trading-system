package services

import (
	"log"
)

type MarketDataService struct {
	Clients []ExchangeClient
}

func NewMarketDataService(clients ...ExchangeClient) *MarketDataService {
	return &MarketDataService{Clients: clients}
}

func (m *MarketDataService) Start(symbol string) {
	for _, client := range m.Clients {
		go func(client ExchangeClient) {
			if err := client.SubscribeToTicker(symbol); err != nil {
				log.Printf("Failed to subscribe: %v", err)
				return
			}
			client.ReadMessages()
		}(client)
	}
}

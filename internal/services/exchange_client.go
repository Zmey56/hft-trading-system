package services

// ExchangeClient Interface for the exchange client
type ExchangeClient interface {
	SubscribeToTicker(symbol string) error
	ReadMessages()
}

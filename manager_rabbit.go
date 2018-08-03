package manager

// RabbitmqConfig ...
type RabbitmqConfig struct {
	Uri          string `json:"uri"`
	Exchange     string `json:"exchange"`
	ExchangeType string `json:"exchange_type"`
}

// NewRabbitmqConfig...
func NewRabbitmqConfig(uri, exchange, exchangeType string) *RabbitmqConfig {
	return &RabbitmqConfig{
		Uri:          uri,
		Exchange:     exchange,
		ExchangeType: exchangeType,
	}
}

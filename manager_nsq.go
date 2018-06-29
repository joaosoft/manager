package manager

// NSQConfig ...
type NSQConfig struct {
	Lookupd      []string `json:"lookupd"`
	Nsqd         []string `json:"nsqd"`
	Topic        string   `json:"topic"`
	Channel      string   `json:"channel"`
	RequeueDelay int64    `json:"requeue_delay"`
	MaxInFlight  int      `json:"max_in_flight"`
	MaxAttempts  uint16   `json:"max_attempts"`
	AutoRespond  bool     `json:"auto_respond"`
}

// NewNSQConfig...
func NewNSQConfig(topic, channel string, lookupd, nsqd []string, requeueDelay int64, maxInFlight int) *NSQConfig {
	return &NSQConfig{
		Topic:        topic,
		Channel:      channel,
		Lookupd:      lookupd,
		Nsqd:         nsqd,
		RequeueDelay: requeueDelay,
		MaxInFlight:  maxInFlight,
	}
}

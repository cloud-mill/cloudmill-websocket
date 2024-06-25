package models

import "time"

type Message struct {
	Id        string                 `json:"id"        mapstructure:"id"`
	SendTo    []string               `json:"send_to"   mapstructure:"send_to"`
	Type      string                 `json:"type"      mapstructure:"type"`
	Timestamp time.Time              `json:"timestamp" mapstructure:"timestamp"`
	Payload   map[string]interface{} `json:"payload"   mapstructure:"payload"`
}

type ProcessedMessage struct {
	Id        string                 `json:"id"        mapstructure:"id"`
	Type      string                 `json:"type"      mapstructure:"type"`
	Timestamp time.Time              `json:"timestamp" mapstructure:"timestamp"`
	Payload   map[string]interface{} `json:"payload"   mapstructure:"payload"`
}

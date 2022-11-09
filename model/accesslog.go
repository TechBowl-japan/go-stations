package model

import "time"

type (
	AccessLog struct {
		Latency   int64     `json:"latency"`
		OS        string    `json:"os"`
		Timestamp time.Time `json:"timestamp"`
		Path      string    `json:"path"`
	}
)

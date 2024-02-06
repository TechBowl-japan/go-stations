package model

// A HealthzResponse expresses health check message.
type HealthzResponse struct {
	Message string `json:"message"`
}

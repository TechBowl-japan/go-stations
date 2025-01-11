package model

// A HealthzResponse expresses health check message.
type HealthzResponse struct{
	Message string `json:"message"` // JSONで返却する際のキー名を指定
}

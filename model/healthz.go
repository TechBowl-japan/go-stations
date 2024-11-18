package model

// A HealthzResponse expresses health check message.
// HealthzResponseはヘルスチェックメッセージを表します。
type HealthzResponse struct {
	Message string `json:"message"`
}

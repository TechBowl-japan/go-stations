package model

import "time"

type (
	// A TODO expresses ...
	//TODOは保存されるTODOのデータ形式を表現します。
	TODO struct {
		ID          int64     `json:"id"`
		Subject     string    `json:"subject"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"` //キャメルケース(ラクダケース：CreatedAt(大文字小文字大文字の山谷で覚える）により、Created_atではなく、CreatedAt
		UpdatedAt   time.Time `json:"updated_at"` //スネークケース（蛇ケース：updated_at（_が地面を這いつくばっている蛇のように覚える）
	}

	// A CreateTODORequest expresses ...
	// CreateTODORequestは利用者からのリクエスト形式
	CreateTODORequest struct {
		Subject     string `json:"subject"`
		Description string `json:"description"`
	}
	// A CreateTODOResponse expresses ...
	// CreateTODOResponseは保存したTODOをレスポンスとして返す
	CreateTODOResponse struct {
		TODO TODO `json:"todo"`
	}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct {
		PrevID int64 `json:"prev_id"`
		Size   int64 `json:"size"`
	}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct {
		TODOs []TODO `json:"todos"`
	}

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct {
		ID          int64  `json:"id"`
		Subject     string `json:"subject"`
		Description string `json:"description"`
	}
	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct {
		TODO TODO `json:"todo"`
	}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct {
		IDs []int64 `json:"ids"`
	}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
	//`json:"os"`などにするのは、JSON形式に変換（シリアライズ）する際に、osというキー名で出力されるべきという指示を持つ。
	//jsonタグがあれば、os、なければOSのフィールド名が出力される
	Requestlog struct {
		Timestamp time.Time `json:"timestamp"`
		Latency   int64     `json:"latency"`
		Path      string    `json:"path"`
		OS        string    `json:"os"`
	}
)

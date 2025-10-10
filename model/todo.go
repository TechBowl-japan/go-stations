package model

import "time"

type (
	// A TODO expresses a TODO item.
	TODO struct{
		ID          int       `json:"id"`          // JSONキーは "id"
		Subject     string    `json:"subject"`     // JSONキーは "subject"
		Description string    `json:"description"` // JSONキーは "description"
		CreatedAt   time.Time `json:"created_at"`  // JSONキーは "created_at"
		UpdatedAt   time.Time `json:"updated_at"`  // JSONキーは "updated_at"
	}

	// A CreateTODORequest expresses the request to create a new TODO item.
	CreateTODORequest struct{
		Subject     string `json:"subject"`     // 必須フィールド
		Description string `json:"description"` // 任意のフィールド
	}
	// A CreateTODOResponse expresses the response after creating a new TODO item.
	CreateTODOResponse struct{
		TODO TODO `json:"todo"` // 作成された TODO
	}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct{}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct{}

	// A UpdateTODORequest expresses the request to update a TODO item.
	UpdateTODORequest struct{
		ID          int    `json:"id"`          // 更新対象の TODO ID（必須）
		Subject     string `json:"subject"`     // 更新後のタイトル（必須）
		Description string `json:"description"` // 更新後の説明（任意）
	}
	// A UpdateTODOResponse expresses the response after updating a TODO item.
	UpdateTODOResponse struct{
		TODO TODO `json:"todo"` // 更新後の TODO
	}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
)

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

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct{}
	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct{}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
)

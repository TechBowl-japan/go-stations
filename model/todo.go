package model

import "time"

type (
	// A TODO expresses a TODO item.
	TODO struct{
		ID          int64       `json:"id"`          // JSONキーは "id"
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
	ReadTODORequest struct{
		PrevID int64 `json:"prev_id"` // 直前のID（以降を読むためのカーソル）
    	Size   int64 `json:"size"`    // 取得件数
	}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct{
		TODOs []TODO `json:"todos"` // 取得したTODOの配列（*TODOでも可）
	}

	// A UpdateTODORequest expresses the request to update a TODO item.
	UpdateTODORequest struct{
		ID          int64    `json:"id"`          // 更新対象の TODO ID（必須）
		Subject     string `json:"subject"`     // 更新後のタイトル（必須）
		Description string `json:"description"` // 更新後の説明（任意）
	}
	// A UpdateTODOResponse expresses the response after updating a TODO item.
	UpdateTODOResponse struct{
		TODO TODO `json:"todo"` // 更新後の TODO
	}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{
		IDs []int64 `json:"ids"` // 削除対象の TODO ID の配列
	}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
)

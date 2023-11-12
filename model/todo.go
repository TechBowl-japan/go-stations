package model

import "time"

type (
	// A TODO expresses ...
	TODO struct{
		ID int64 `json:"id"`
		Subject string`json:"subject"`
		Description string`json:"description"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	// A CreateTODORequest expresses ...
	CreateTODORequest struct{
		Subject string `json:"subject"`
		Description string `json:"description"`
	}
	// A CreateTODOResponse expresses ...
	CreateTODOResponse struct{
		TODO TODO `json:"todo"`
		// TODOList []TODO `json:"todo"`　←これがダメな理由を教えて
		// responseのtype:objectと指定されているから配列ではダメということなのか？
		// propertiesでtodoが指定されているからタグ指定もtodoにしないといけないのか？
	}

	// A ReadTODORequest expresses ...
	ReadTODORequest struct{
	}
	// A ReadTODOResponse expresses ...
	ReadTODOResponse struct{}

	// A UpdateTODORequest expresses ...
	UpdateTODORequest struct{
		ID int64 `json:"id"`
		Subject string `json:"subject"`
		Description string `json:"description"`
	}
	// A UpdateTODOResponse expresses ...
	UpdateTODOResponse struct{
		TODO TODO `json:"todo"`
	}

	// A DeleteTODORequest expresses ...
	DeleteTODORequest struct{}
	// A DeleteTODOResponse expresses ...
	DeleteTODOResponse struct{}
)

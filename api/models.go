package api

//go:generate easyjson -all models.go

type EmptySuccessResponse struct {
	Status string `json:"status"`
}

func NewEmptySuccessResponse() EmptySuccessResponse {
	return EmptySuccessResponse{
		Status: "success",
	}
}

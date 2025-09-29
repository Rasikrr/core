package api

//go:generate easyjson -all models.go

type EmptySuccessResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status"`
}

type SuccessResponse struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
}

func NewEmptySuccessResponse() EmptySuccessResponse {
	return EmptySuccessResponse{
		Status: "success",
	}
}

func NewSuccessResponse(data interface{}) SuccessResponse {
	return SuccessResponse{
		Data:   data,
		Status: "success",
	}
}

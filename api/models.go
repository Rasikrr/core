package api

//go:generate easyjson -all models.go

type EmptySuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func NewEmptySuccessResponse(message ...string) EmptySuccessResponse {
	resp := EmptySuccessResponse{
		Message: "success",
	}
	if len(message) > 0 {
		resp.Message = message[0]
	}
	return resp
}

func NewSuccessResponse(data interface{}, message ...string) SuccessResponse {
	resp := SuccessResponse{
		Data:    data,
		Message: "success",
	}
	if len(message) > 0 {
		resp.Message = message[0]
	}
	return resp
}

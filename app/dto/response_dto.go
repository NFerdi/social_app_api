package dto

type SuccessResponseDto struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponseDto struct {
	Type   string      `json:"type"`
	Errors interface{} `json:"error"`
}

package domain

type ErrorResponse struct {
	Error struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Detail  interface{} `json:"detail"`
	} `json:"error"`
}

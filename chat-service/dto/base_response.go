package dto

type BaseResponse struct {
	Status int         `json:"status"`
	Result interface{} `json:"result"`
	Error  interface{} `json:"error"`
}

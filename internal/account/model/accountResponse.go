package model

type AccountResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error_message"`
}

package model

type EventResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error_message"`
}

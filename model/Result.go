package model

type Result struct {
	StatusCode int         `json:"status_code"`
	Data       interface{} `json:"data"`
}

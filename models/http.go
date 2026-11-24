package models

type Response[T any] struct {
	Success bool   `json:"success"`
	IsError bool   `json:"isError"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	Data    *T     `json:"data,omitempty"`
}

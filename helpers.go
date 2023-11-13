package main

import "fmt"

type ApiError struct {
	Error string `json:"error"`
}

func NewApiError(err error) ApiError {
	return ApiError{Error: fmt.Sprintf("%s", err)}
}

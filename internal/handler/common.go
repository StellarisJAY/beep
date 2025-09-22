package handler

import "net/http"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ok() *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: "success",
	}
}

func okWithData(data interface{}) *Response {
	return &Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
}

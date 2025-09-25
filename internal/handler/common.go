package handler

import "net/http"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int         `json:"total"`
}

func (r *Response) withTotal(total int) *Response {
	r.Total = total
	return r
}

func (r *Response) withData(result interface{}) *Response {
	r.Data = result
	return r
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

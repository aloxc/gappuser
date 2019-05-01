package io

type Response struct {
	Code    int8        `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
type ResponseType interface {
	code() int8
	message() string
}

type RuntimeResponse struct {
}

func (this *RuntimeResponse) code() int8 {
	return this.code()
}
func (this *RuntimeResponse) message() string {
	return this.message()
}

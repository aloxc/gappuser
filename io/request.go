package io

type Request struct {
	Method string                 `json:"method"`
	Params map[string]interface{} `json:"params"`
	Head   RequestHead            `json:"head"`
}
type RequestHead struct {
	UserIp     string `json:"user_ip"`
	UserPort   uint16 `json:"user_port"`
	UserCookie string `json:"user_cookie"`
}

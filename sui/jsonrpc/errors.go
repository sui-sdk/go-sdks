package jsonrpc

import "fmt"

type HTTPStatusError struct {
	StatusCode int
	Body       string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("sui http status error: %d", e.StatusCode)
}

type HTTPTransportError struct {
	Cause error
}

func (e *HTTPTransportError) Error() string {
	if e.Cause == nil {
		return "sui http transport error"
	}
	return "sui http transport error: " + e.Cause.Error()
}

type JsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *JsonRPCError) Error() string {
	return fmt.Sprintf("json-rpc error %d: %s", e.Code, e.Message)
}

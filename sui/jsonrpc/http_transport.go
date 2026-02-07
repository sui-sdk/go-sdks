package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type HTTPHeaders map[string]string

type TransportRequest struct {
	Method string
	Params []any
	Ctx    context.Context
}

type Transport interface {
	Request(req TransportRequest, out any) error
}

type HTTPTransportOptions struct {
	URL     string
	Headers HTTPHeaders
	Client  *http.Client
}

type HTTPTransport struct {
	url     string
	headers HTTPHeaders
	client  *http.Client
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
}

type rpcResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Result  any           `json:"result"`
	Error   *JsonRPCError `json:"error,omitempty"`
}

func NewHTTPTransport(opts HTTPTransportOptions) *HTTPTransport {
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	if opts.Headers == nil {
		opts.Headers = HTTPHeaders{}
	}
	return &HTTPTransport{url: opts.URL, headers: opts.Headers, client: client}
}

func (t *HTTPTransport) Request(req TransportRequest, out any) error {
	payload, err := json.Marshal(rpcRequest{JSONRPC: "2.0", ID: 1, Method: req.Method, Params: req.Params})
	if err != nil {
		return &HTTPTransportError{Cause: err}
	}

	httpReq, err := http.NewRequestWithContext(req.Ctx, http.MethodPost, t.url, bytes.NewReader(payload))
	if err != nil {
		return &HTTPTransportError{Cause: err}
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range t.headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return &HTTPTransportError{Cause: err}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HTTPStatusError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var rpcResp rpcResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return &HTTPTransportError{Cause: err}
	}
	if rpcResp.Error != nil {
		return rpcResp.Error
	}
	if out == nil {
		return nil
	}

	buf, err := json.Marshal(rpcResp.Result)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, out)
}

package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type QueryOptions struct {
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables,omitempty"`
	OperationName string         `json:"operationName,omitempty"`
	Extensions    map[string]any `json:"extensions,omitempty"`
}

type QueryResult struct {
	Data       any              `json:"data,omitempty"`
	Errors     []ResponseError  `json:"errors,omitempty"`
	Extensions map[string]any   `json:"extensions,omitempty"`
}

type ResponseError struct {
	Message   string           `json:"message"`
	Locations []map[string]int `json:"locations,omitempty"`
	Path      []any            `json:"path,omitempty"`
}

type RequestError struct {
	StatusCode int
	Status     string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("graphql request failed: %s (%d)", e.Status, e.StatusCode)
}

type ClientOptions struct {
	URL     string
	Network string
	Headers map[string]string
	Fetch   *http.Client
	Queries map[string]string
}

type Client struct {
	url     string
	network string
	headers map[string]string
	client  *http.Client
	queries map[string]string
}

func NewClient(opts ClientOptions) *Client {
	client := opts.Fetch
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	if opts.Headers == nil {
		opts.Headers = map[string]string{}
	}
	if opts.Queries == nil {
		opts.Queries = map[string]string{}
	}
	return &Client{url: opts.URL, network: opts.Network, headers: opts.Headers, client: client, queries: opts.Queries}
}

func (c *Client) Network() string { return c.network }

func (c *Client) Query(ctx context.Context, opts QueryOptions) (*QueryResult, error) {
	payload, _ := json.Marshal(opts)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &RequestError{StatusCode: resp.StatusCode, Status: resp.Status}
	}
	var out QueryResult
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) Execute(ctx context.Context, queryName string, variables map[string]any, operationName string, extensions map[string]any) (*QueryResult, error) {
	q, ok := c.queries[queryName]
	if !ok {
		return nil, fmt.Errorf("unknown query: %s", queryName)
	}
	return c.Query(ctx, QueryOptions{
		Query:         q,
		Variables:     variables,
		OperationName: operationName,
		Extensions:    extensions,
	})
}

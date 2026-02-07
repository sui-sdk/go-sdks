package walrus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type storageNodeClientOptionsAlias struct {
	Client  *http.Client
	Timeout time.Duration
	OnError func(error)
}

type RequestOptions struct {
	Ctx     context.Context
	NodeURL string
	Headers map[string]string
	Timeout time.Duration
}

type StorageNodeClient struct {
	client  *http.Client
	timeout time.Duration
	onError func(error)
}

func NewStorageNodeClient(opts *storageNodeClientOptionsAlias) *StorageNodeClient {
	if opts == nil {
		opts = &storageNodeClientOptionsAlias{}
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	client := opts.Client
	if client == nil {
		client = &http.Client{Timeout: timeout}
	}
	return &StorageNodeClient{client: client, timeout: timeout, onError: opts.OnError}
}

func (c *StorageNodeClient) GetBlobMetadata(input map[string]string, opts RequestOptions) ([]byte, error) {
	blobID := input["blobId"]
	resp, err := c.request(opts, http.MethodGet, fmt.Sprintf("/v1/blobs/%s/metadata", blobID), nil, map[string]string{"Accept": "application/octet-stream"})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *StorageNodeClient) GetBlobStatus(input map[string]string, opts RequestOptions) (map[string]any, error) {
	blobID := input["blobId"]
	resp, err := c.request(opts, http.MethodGet, fmt.Sprintf("/v1/blobs/%s/status", blobID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *StorageNodeClient) StoreBlobMetadata(input map[string]any, opts RequestOptions) (map[string]any, error) {
	blobID, _ := input["blobId"].(string)
	metadata, _ := input["metadata"].([]byte)
	resp, err := c.request(opts, http.MethodPut, fmt.Sprintf("/v1/blobs/%s/metadata", blobID), metadata, map[string]string{"Content-Type": "application/octet-stream"})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *StorageNodeClient) GetSliver(input map[string]any, opts RequestOptions) ([]byte, error) {
	blobID, _ := input["blobId"].(string)
	pairIndex, _ := input["sliverPairIndex"].(int)
	sliverType, _ := input["sliverType"].(string)
	resp, err := c.request(opts, http.MethodGet, fmt.Sprintf("/v1/blobs/%s/slivers/%d/%s", blobID, pairIndex, sliverType), nil, map[string]string{"Accept": "application/octet-stream"})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *StorageNodeClient) StoreSliver(input map[string]any, opts RequestOptions) (map[string]any, error) {
	blobID, _ := input["blobId"].(string)
	pairIndex, _ := input["sliverPairIndex"].(int)
	sliverType, _ := input["sliverType"].(string)
	sliver, _ := input["sliver"].([]byte)
	resp, err := c.request(opts, http.MethodPut, fmt.Sprintf("/v1/blobs/%s/slivers/%d/%s", blobID, pairIndex, sliverType), sliver, map[string]string{"Content-Type": "application/octet-stream"})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *StorageNodeClient) GetDeletableBlobConfirmation(input map[string]string, opts RequestOptions) (map[string]any, error) {
	blobID := input["blobId"]
	objectID := input["objectId"]
	resp, err := c.request(opts, http.MethodGet, fmt.Sprintf("/v1/blobs/%s/confirmation/deletable/%s", blobID, objectID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *StorageNodeClient) GetPermanentBlobConfirmation(input map[string]string, opts RequestOptions) (map[string]any, error) {
	blobID := input["blobId"]
	resp, err := c.request(opts, http.MethodGet, fmt.Sprintf("/v1/blobs/%s/confirmation/permanent", blobID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil
}

func (c *StorageNodeClient) request(opts RequestOptions, method, path string, body []byte, headers map[string]string) (*http.Response, error) {
	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = c.timeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	url := opts.NodeURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range opts.Headers {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		if c.onError != nil {
			c.onError(err)
		}
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		defer resp.Body.Close()
		bs, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("storage node api error: %d %s", resp.StatusCode, string(bs))
		if c.onError != nil {
			c.onError(err)
		}
		return nil, err
	}
	return resp, nil
}

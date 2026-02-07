package seal

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ServerType string

const (
	ServerTypeIndependent ServerType = "Independent"
	ServerTypeCommittee   ServerType = "Committee"
)

type KeyType int

const (
	KeyTypeBonehFranklinBLS12381 KeyType = iota
)

var ServerVersionRequirement = NewVersion("0.4.1")

type KeyServer struct {
	ObjectID   string
	Name       string
	URL        string
	Weight     int
	PK         []byte
	ServerType ServerType
}

type DerivedKey struct {
	ObjectID string
	ID       string
	Key      []byte
}

func RetrieveKeyServers(configs map[string]KeyServerConfig) ([]KeyServer, error) {
	out := make([]KeyServer, 0, len(configs))
	for id, cfg := range configs {
		pk, err := HexToBytes(cfg.PublicKeyHex)
		if err != nil && cfg.PublicKeyHex != "" {
			return nil, err
		}
		out = append(out, KeyServer{ObjectID: id, URL: cfg.URL, Weight: cfg.Weight, PK: pk, ServerType: ServerTypeIndependent})
	}
	if len(out) == 0 {
		return nil, (&InvalidKeyServerError{UserError{SealError{Msg: "no key servers found"}}})
	}
	return out, nil
}

func VerifyKeyServer(server KeyServer, timeout time.Duration, apiKeyName, apiKey string) bool {
	if server.URL == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, server.URL+"/service", nil)
	if apiKeyName != "" && apiKey != "" {
		req.Header.Set(apiKeyName, apiKey)
	}
	resp, err := (&http.Client{Timeout: timeout}).Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

func FetchKeysForAllIDs(ctx context.Context, server KeyServer, ids []string, txBytes []byte, session *SessionKey, timeout time.Duration, apiKeyName, apiKey string) ([]DerivedKey, error) {
	params, err := session.CreateRequestParams(txBytes)
	if err != nil {
		return nil, err
	}
	body := map[string]any{"ids": ids, "request": params}
	bs, _ := json.Marshal(body)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, server.URL+"/v1/keys", bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKeyName != "" && apiKey != "" {
		req.Header.Set(apiKeyName, apiKey)
	}
	resp, err := (&http.Client{Timeout: timeout}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("fetch keys failed: %s", resp.Status)
	}
	var out struct {
		Keys map[string]string `json:"keys"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	result := make([]DerivedKey, 0, len(out.Keys))
	for id, keyB64 := range out.Keys {
		k, err := base64.StdEncoding.DecodeString(keyB64)
		if err != nil {
			continue
		}
		result = append(result, DerivedKey{ObjectID: server.ObjectID, ID: id, Key: k})
	}
	return result, nil
}

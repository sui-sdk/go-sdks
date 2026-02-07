package faucet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type FaucetRateLimitError struct{}

func (e *FaucetRateLimitError) Error() string {
	return "too many requests sent to faucet, retry later"
}

type CoinInfo struct {
	Amount           int    `json:"amount"`
	ID               string `json:"id"`
	TransferTxDigest string `json:"transferTxDigest"`
}

type failureBody struct {
	Failure struct {
		Internal string `json:"internal"`
	} `json:"Failure"`
}

type ResponseV2 struct {
	Status    any        `json:"status"`
	CoinsSent []CoinInfo `json:"coins_sent"`
}

func RequestSuiFromFaucetV2(ctx context.Context, host, recipient string, headers map[string]string) (*ResponseV2, error) {
	body, _ := json.Marshal(map[string]any{
		"FixedAmountRequest": map[string]any{
			"recipient": recipient,
		},
	})
	url := host + "/v2/gas"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, &FaucetRateLimitError{}
	}
	var out ResponseV2
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("parse faucet response failed: %w", err)
	}
	if s, ok := out.Status.(string); ok && s == "Success" {
		return &out, nil
	}
	statusBytes, _ := json.Marshal(out.Status)
	var fb failureBody
	if json.Unmarshal(statusBytes, &fb) == nil && fb.Failure.Internal != "" {
		return nil, fmt.Errorf("faucet request failed: %s", fb.Failure.Internal)
	}
	return nil, fmt.Errorf("faucet request failed: %s", string(statusBytes))
}

func GetFaucetHost(network string) (string, error) {
	switch network {
	case "testnet":
		return "https://faucet.testnet.sui.io", nil
	case "devnet":
		return "https://faucet.devnet.sui.io", nil
	case "localnet":
		return "http://127.0.0.1:9123", nil
	default:
		return "", fmt.Errorf("unknown network: %s", network)
	}
}

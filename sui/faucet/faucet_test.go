package faucet

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestSuiFromFaucetV2Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"Success","coins_sent":[{"amount":1000,"id":"0x1","transferTxDigest":"abc"}]}`))
	}))
	defer srv.Close()

	res, err := RequestSuiFromFaucetV2(context.Background(), srv.URL, "0x1", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if len(res.CoinsSent) != 1 {
		t.Fatalf("expected coins")
	}
}

func TestRequestSuiFromFaucetV2RateLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer srv.Close()

	_, err := RequestSuiFromFaucetV2(context.Background(), srv.URL, "0x1", nil)
	if err == nil {
		t.Fatalf("expected rate limit error")
	}
	if _, ok := err.(*FaucetRateLimitError); !ok {
		t.Fatalf("expected FaucetRateLimitError got %T", err)
	}
}

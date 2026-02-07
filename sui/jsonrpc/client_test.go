package jsonrpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientGetBalanceAndCall(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		method := req["method"].(string)
		w.Header().Set("Content-Type", "application/json")
		switch method {
		case "suix_getBalance":
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "result": map[string]any{"coinType": "0x2::sui::SUI", "totalBalance": "123"}})
		case "rpc.discover":
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "result": map[string]any{"info": map[string]any{"version": "1.0.0"}}})
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "result": map[string]any{"ok": true}})
		}
	}))
	defer srv.Close()

	c, err := NewClient(ClientOptions{Network: "testnet", URL: srv.URL})
	if err != nil {
		t.Fatalf("new client failed: %v", err)
	}
	v, err := c.GetRPCAPIVersion(context.Background())
	if err != nil || v != "1.0.0" {
		t.Fatalf("unexpected version: %q err=%v", v, err)
	}
	bal, err := c.GetBalance(context.Background(), "0x0000000000000000000000000000000000000000000000000000000000000001", "")
	if err != nil {
		t.Fatalf("get balance failed: %v", err)
	}
	if bal["totalBalance"].(string) != "123" {
		t.Fatalf("unexpected balance")
	}
}

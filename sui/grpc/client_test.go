package grpc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	jsonrpc "github.com/sui-sdks/go-sdks/sui/jsonrpc"
)

func TestGrpcClientCoreMethodsWithJSONRPCTransport(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		method := req["method"].(string)
		w.Header().Set("Content-Type", "application/json")
		switch method {
		case "suix_getReferenceGasPrice":
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "result": "1000"})
		case "suix_getBalance":
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "result": map[string]any{"totalBalance": "1"}})
		case "sui_getObject":
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "result": map[string]any{"objectId": "0x1"}})
		case "sui_multiGetObjects":
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "result": []map[string]any{{"objectId": "0x1"}}})
		default:
			_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": 1, "result": map[string]any{}})
		}
	}))
	defer srv.Close()

	rpc, err := jsonrpc.NewClient(jsonrpc.ClientOptions{Network: "testnet", URL: srv.URL})
	if err != nil {
		t.Fatalf("new jsonrpc client failed: %v", err)
	}

	client, err := NewClient(ClientOptions{Network: "testnet", RPC: rpc})
	if err != nil {
		t.Fatalf("new client failed: %v", err)
	}
	if _, err := client.GetReferenceGasPrice(context.Background()); err != nil {
		t.Fatalf("get reference gas price failed: %v", err)
	}
	if _, err := client.GetBalance(context.Background(), "0x0000000000000000000000000000000000000000000000000000000000000001", ""); err != nil {
		t.Fatalf("get balance failed: %v", err)
	}
	objectID := "0x0000000000000000000000000000000000000000000000000000000000000001"
	if _, err := client.GetObject(context.Background(), objectID, nil); err != nil {
		t.Fatalf("get object failed: %v", err)
	}
	if _, err := client.GetObjects(context.Background(), []string{objectID}, nil); err != nil {
		t.Fatalf("get objects failed: %v", err)
	}
}

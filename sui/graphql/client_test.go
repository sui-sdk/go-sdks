package graphql

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGraphQLQueryAndExecute(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		_ = json.NewEncoder(w).Encode(map[string]any{"data": map[string]any{"echo": req["query"]}})
	}))
	defer srv.Close()

	c := NewClient(ClientOptions{URL: srv.URL, Network: "testnet", Queries: map[string]string{"q1": "query { ping }"}})
	res, err := c.Query(context.Background(), QueryOptions{Query: "query { ping }"})
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	if res.Data == nil {
		t.Fatalf("expected data")
	}
	res2, err := c.Execute(context.Background(), "q1", nil, "", nil)
	if err != nil || res2.Data == nil {
		t.Fatalf("execute failed: %v", err)
	}
}

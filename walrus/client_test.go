package walrus

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWalrusClientReadAndWrite(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/blobs/b1/status":
			_, _ = w.Write([]byte(`{"success":{"data":{"deletable":{"dummy":true}}}}`))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/blobs/b1/metadata":
			_, _ = w.Write([]byte{1, 2, 3})
		case r.Method == http.MethodPut && r.URL.Path == "/v1/blobs/b1/metadata":
			_, _ = w.Write([]byte(`{"ok":true}`))
		case r.Method == http.MethodPut && r.URL.Path == "/v1/blobs/b1/slivers/0/primary":
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	c, err := NewClient(Config{Network: "testnet"})
	if err != nil {
		t.Fatalf("new client failed: %v", err)
	}
	blob, err := c.ReadBlob(context.Background(), srv.URL, ReadBlobOptions{BlobID: "b1"})
	if err != nil {
		t.Fatalf("read blob failed: %v", err)
	}
	if len(blob) != 3 {
		t.Fatalf("unexpected blob length")
	}
	if _, err := c.WriteBlobMetadata(context.Background(), WriteBlobMetadataOptions{NodeURL: srv.URL, BlobID: "b1", Metadata: []byte{1}}); err != nil {
		t.Fatalf("write metadata failed: %v", err)
	}
	if _, err := c.WriteSliver(context.Background(), WriteSliverOptions{NodeURL: srv.URL, BlobID: "b1", SliverPairIndex: 0, SliverType: "primary", Sliver: []byte{1}}); err != nil {
		t.Fatalf("write sliver failed: %v", err)
	}
}

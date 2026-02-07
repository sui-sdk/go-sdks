package transactions

import (
	"context"
	"encoding/base64"
	"testing"

	edkp "github.com/sui-sdks/go-sdks/sui/keypairs/ed25519"
)

type mockCore struct{}

func (m mockCore) Call(ctx context.Context, method string, params []any, out any) error {
	switch method {
	case "suix_getReferenceGasPrice":
		if p, ok := out.(*string); ok {
			*p = "1000"
		}
	case "sui_executeTransactionBlock":
		if p, ok := out.(*map[string]any); ok {
			*p = map[string]any{"digest": "abc"}
		}
	case "sui_getTransactionBlock":
		if p, ok := out.(*map[string]any); ok {
			*p = map[string]any{"digest": "abc"}
		}
	default:
		if p, ok := out.(*map[string]any); ok {
			*p = map[string]any{}
		}
	}
	_ = params
	_ = base64.StdEncoding
	return nil
}

func TestResolveTransactionPlugin(t *testing.T) {
	tx := NewTransaction()
	tx.data.Inputs = append(tx.data.Inputs, CallArg{"$kind": "UnresolvedPure", "UnresolvedPure": map[string]any{"value": 10}})
	err := ResolveTransactionPlugin(&tx.data, BuildTransactionOptions{Client: mockCore{}}, func() error { return nil })
	if err != nil {
		t.Fatalf("resolve plugin failed: %v", err)
	}
}

func TestCachingSerialParallelExecutors(t *testing.T) {
	signer, err := edkp.Generate()
	if err != nil {
		t.Fatalf("generate signer failed: %v", err)
	}
	tx := NewTransaction()
	tx.SetSender("0x1")
	tx.SetGasBudget(1000)
	tx.PureBytes([]byte("x"))

	cacheExec := NewCachingTransactionExecutor(mockCore{}, nil)
	b, err := cacheExec.BuildTransaction(tx, BuildTransactionOptions{Client: mockCore{}})
	if err != nil || len(b) == 0 {
		t.Fatalf("build transaction failed: %v", err)
	}
	if _, err := cacheExec.ExecuteTransaction(ExecuteTransactionOptions{Transaction: b, Signatures: []string{"sig"}, Include: nil}); err != nil {
		t.Fatalf("execute transaction failed: %v", err)
	}

	serial := NewSerialTransactionExecutor(SerialTransactionExecutorOptions{Client: mockCore{}, Signer: signer})
	if _, err := serial.ExecuteTransaction(tx, nil, nil); err != nil {
		t.Fatalf("serial execute failed: %v", err)
	}

	parallel := NewParallelTransactionExecutor(ParallelTransactionExecutorOptions{Client: mockCore{}, Signer: signer})
	if _, err := parallel.ExecuteTransaction(tx, nil, nil); err != nil {
		t.Fatalf("parallel execute failed: %v", err)
	}
}

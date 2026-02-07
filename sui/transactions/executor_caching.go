package transactions

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
)

type ExecuteTransactionOptions struct {
	Transaction []byte
	Signatures  []string
	Include     map[string]any
}

type ExecuteCore interface {
	Call(ctx context.Context, method string, params []any, out any) error
}

type CachingTransactionExecutor struct {
	client     ExecuteCore
	cache      *ObjectCache
	lastDigest string
}

func NewCachingTransactionExecutor(client ExecuteCore, cache *ObjectCache) *CachingTransactionExecutor {
	if cache == nil {
		cache = NewObjectCache()
	}
	return &CachingTransactionExecutor{client: client, cache: cache}
}

func (e *CachingTransactionExecutor) Reset() error {
	e.cache.ClearOwnedObjects()
	e.cache.ClearCustom()
	return e.WaitForLastTransaction()
}

func (e *CachingTransactionExecutor) BuildTransaction(tx *Transaction, options BuildTransactionOptions) ([]byte, error) {
	data := tx.GetData()
	if err := ResolveTransactionPlugin(&data, options, func() error { return nil }); err != nil {
		return nil, err
	}
	tx.data = data
	return tx.Build()
}

func (e *CachingTransactionExecutor) ExecuteTransaction(opts ExecuteTransactionOptions) (map[string]any, error) {
	if len(opts.Signatures) == 0 {
		return nil, errors.New("at least one signature is required")
	}
	txB64 := base64.StdEncoding.EncodeToString(opts.Transaction)
	var out map[string]any
	err := e.client.Call(context.Background(), "sui_executeTransactionBlock", []any{txB64, opts.Signatures, opts.Include, "WaitForLocalExecution"}, &out)
	if err != nil {
		return nil, err
	}
	if digest, ok := out["digest"].(string); ok {
		e.lastDigest = digest
	}
	return out, nil
}

func (e *CachingTransactionExecutor) SignAndExecuteTransaction(tx *Transaction, signer interface {
	ToSuiAddress() string
	SignTransaction([]byte) (cryptography.SignatureWithBytes, error)
}, include map[string]any) (map[string]any, error) {
	tx.SetSenderIfNotSet(signer.ToSuiAddress())
	bytes, err := e.BuildTransaction(tx, BuildTransactionOptions{Client: e.client})
	if err != nil {
		return nil, err
	}
	sig, err := signer.SignTransaction(bytes)
	if err != nil {
		return nil, err
	}
	return e.ExecuteTransaction(ExecuteTransactionOptions{Transaction: bytes, Signatures: []string{sig.Signature}, Include: include})
}

func (e *CachingTransactionExecutor) WaitForLastTransaction() error {
	if e.lastDigest == "" {
		return nil
	}
	var out map[string]any
	err := e.client.Call(context.Background(), "sui_getTransactionBlock", []any{e.lastDigest, map[string]any{}}, &out)
	if err == nil {
		e.lastDigest = ""
	}
	return err
}

func (e *CachingTransactionExecutor) ApplyEffects(effects map[string]any) {
	_ = effects
}

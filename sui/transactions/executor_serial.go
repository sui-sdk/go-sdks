package transactions

import (
	"errors"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
)

type SerialTransactionExecutorOptions struct {
	Client           ExecuteCore
	Signer           interface {
		ToSuiAddress() string
		SignTransaction([]byte) (cryptography.SignatureWithBytes, error)
	}
	DefaultGasBudget int64
	GasMode          string
	Cache            *ObjectCache
}

type SerialTransactionExecutor struct {
	queue            SerialQueue
	signer           interface {
		ToSuiAddress() string
		SignTransaction([]byte) (cryptography.SignatureWithBytes, error)
	}
	cacheExec        *CachingTransactionExecutor
	defaultGasBudget int64
	gasMode          string
}

func NewSerialTransactionExecutor(opts SerialTransactionExecutorOptions) *SerialTransactionExecutor {
	budget := opts.DefaultGasBudget
	if budget <= 0 {
		budget = 50_000_000
	}
	mode := opts.GasMode
	if mode == "" {
		mode = "coins"
	}
	return &SerialTransactionExecutor{
		signer:           opts.Signer,
		cacheExec:        NewCachingTransactionExecutor(opts.Client, opts.Cache),
		defaultGasBudget: budget,
		gasMode:          mode,
	}
}

func (e *SerialTransactionExecutor) ApplyEffects(effects map[string]any) {
	e.cacheExec.ApplyEffects(effects)
}

func (e *SerialTransactionExecutor) BuildTransaction(tx *Transaction) ([]byte, error) {
	var out []byte
	err := e.queue.RunTask(func() error {
		copyTx := *tx
		copyTx.SetGasBudgetIfNotSet(e.defaultGasBudget)
		copyTx.SetSenderIfNotSet(e.signer.ToSuiAddress())
		built, err := e.cacheExec.BuildTransaction(&copyTx, BuildTransactionOptions{Client: e.cacheExec.client, OnlyTransactionKind: false})
		if err != nil {
			return err
		}
		out = built
		return nil
	})
	return out, err
}

func (e *SerialTransactionExecutor) ExecuteTransaction(txOrBytes any, include map[string]any, additionalSignatures []string) (map[string]any, error) {
	if len(additionalSignatures) == 0 {
		additionalSignatures = []string{}
	}
	var out map[string]any
	err := e.queue.RunTask(func() error {
		var bytes []byte
		switch v := txOrBytes.(type) {
		case *Transaction:
			copyTx := *v
			copyTx.SetGasBudgetIfNotSet(e.defaultGasBudget)
			copyTx.SetSenderIfNotSet(e.signer.ToSuiAddress())
			b, err := e.cacheExec.BuildTransaction(&copyTx, BuildTransactionOptions{Client: e.cacheExec.client, OnlyTransactionKind: false})
			if err != nil {
				return err
			}
			bytes = b
		case []byte:
			bytes = v
		default:
			return errors.New("unsupported transaction type")
		}
		sig, err := e.signer.SignTransaction(bytes)
		if err != nil {
			return err
		}
		sigs := append([]string{sig.Signature}, additionalSignatures...)
		res, err := e.cacheExec.ExecuteTransaction(ExecuteTransactionOptions{Transaction: bytes, Signatures: sigs, Include: include})
		if err != nil {
			_ = e.cacheExec.Reset()
			return err
		}
		out = res
		return nil
	})
	return out, err
}

func (e *SerialTransactionExecutor) ResetCache() error {
	return e.cacheExec.Reset()
}

func (e *SerialTransactionExecutor) WaitForLastTransaction() error {
	return e.cacheExec.WaitForLastTransaction()
}

package transactions

import (
	"sync"

	"github.com/sui-sdks/go-sdks/sui/cryptography"
)

type ParallelTransactionExecutorOptions struct {
	Client           ExecuteCore
	Signer           interface {
		ToSuiAddress() string
		SignTransaction([]byte) (cryptography.SignatureWithBytes, error)
	}
	DefaultGasBudget int64
	MaxPoolSize      int
	Cache            *ObjectCache
}

type ParallelTransactionExecutor struct {
	signer           interface {
		ToSuiAddress() string
		SignTransaction([]byte) (cryptography.SignatureWithBytes, error)
	}
	cacheExec        *CachingTransactionExecutor
	buildQueue       SerialQueue
	execQueue        *ParallelQueue
	defaultGasBudget int64
	objectQueues     map[string][]func()
	objectMu         sync.Mutex
}

func NewParallelTransactionExecutor(opts ParallelTransactionExecutorOptions) *ParallelTransactionExecutor {
	budget := opts.DefaultGasBudget
	if budget <= 0 {
		budget = 50_000_000
	}
	maxPool := opts.MaxPoolSize
	if maxPool <= 0 {
		maxPool = 50
	}
	return &ParallelTransactionExecutor{
		signer:           opts.Signer,
		cacheExec:        NewCachingTransactionExecutor(opts.Client, opts.Cache),
		execQueue:        NewParallelQueue(maxPool),
		defaultGasBudget: budget,
		objectQueues:     map[string][]func(){},
	}
}

func (e *ParallelTransactionExecutor) ResetCache() error {
	return e.cacheExec.Reset()
}

func (e *ParallelTransactionExecutor) WaitForLastTransaction() error {
	return e.cacheExec.WaitForLastTransaction()
}

func (e *ParallelTransactionExecutor) ExecuteTransaction(tx *Transaction, include map[string]any, additionalSignatures []string) (map[string]any, error) {
	usedObjects := getUsedObjects(tx)
	for _, objectID := range usedObjects {
		e.lockObject(objectID)
		defer e.unlockObject(objectID)
	}

	var out map[string]any
	err := e.execQueue.RunTask(func() error {
		var built []byte
		if err := e.buildQueue.RunTask(func() error {
			copyTx := *tx
			copyTx.SetSenderIfNotSet(e.signer.ToSuiAddress())
			copyTx.SetGasBudgetIfNotSet(e.defaultGasBudget)
			b, err := e.cacheExec.BuildTransaction(&copyTx, BuildTransactionOptions{Client: e.cacheExec.client})
			if err != nil {
				return err
			}
			built = b
			return nil
		}); err != nil {
			return err
		}
		sig, err := e.signer.SignTransaction(built)
		if err != nil {
			return err
		}
		sigs := append([]string{sig.Signature}, additionalSignatures...)
		res, err := e.cacheExec.ExecuteTransaction(ExecuteTransactionOptions{Transaction: built, Signatures: sigs, Include: include})
		if err != nil {
			_ = e.cacheExec.Reset()
			return err
		}
		out = res
		return nil
	})
	return out, err
}

func getUsedObjects(tx *Transaction) []string {
	seen := map[string]struct{}{}
	out := []string{}
	for _, input := range tx.data.Inputs {
		if obj, ok := input["Object"].(map[string]any); ok {
			if imm, ok := obj["ImmOrOwnedObject"].(map[string]any); ok {
				if id, ok := imm["objectId"].(string); ok {
					if _, exists := seen[id]; !exists {
						seen[id] = struct{}{}
						out = append(out, id)
					}
				}
			}
		}
	}
	return out
}

func (e *ParallelTransactionExecutor) lockObject(id string) {
	e.objectMu.Lock()
	defer e.objectMu.Unlock()
	e.objectQueues[id] = append(e.objectQueues[id], func() {})
}

func (e *ParallelTransactionExecutor) unlockObject(id string) {
	e.objectMu.Lock()
	defer e.objectMu.Unlock()
	q := e.objectQueues[id]
	if len(q) <= 1 {
		delete(e.objectQueues, id)
		return
	}
	e.objectQueues[id] = q[1:]
}

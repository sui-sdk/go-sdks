package transactions

import (
	"context"
	"errors"
	"fmt"
)

const (
	GasSafeOverhead = int64(1000)
	MaxGas          = int64(50_000_000_000)
)

func CoreClientResolveTransaction(transactionData *TransactionData, options BuildTransactionOptions) error {
	if options.Client == nil {
		return errors.New("missing core client")
	}
	if err := normalizeInputs(transactionData); err != nil {
		return err
	}
	if err := resolveObjectReferences(transactionData, options.Client); err != nil {
		return err
	}
	if !options.OnlyTransactionKind {
		if err := setGasData(transactionData, options.Client); err != nil {
			return err
		}
	}
	return nil
}

func normalizeInputs(transactionData *TransactionData) error {
	for i, input := range transactionData.Inputs {
		if input["$kind"] == "UnresolvedPure" {
			transactionData.Inputs[i] = Inputs.Pure([]byte(toString(input["UnresolvedPure"])))
		}
	}
	return nil
}

func resolveObjectReferences(transactionData *TransactionData, client CoreClient) error {
	for i, input := range transactionData.Inputs {
		if input["$kind"] != "UnresolvedObject" {
			continue
		}
		payload, _ := input["UnresolvedObject"].(map[string]any)
		objectID, _ := payload["objectId"].(string)
		if objectID == "" {
			return fmt.Errorf("invalid unresolved object at input %d", i)
		}
		transactionData.Inputs[i] = map[string]any{
			"$kind": "Object",
			"Object": map[string]any{
				"$kind": "ImmOrOwnedObject",
				"ImmOrOwnedObject": map[string]any{
					"objectId": objectID,
					"version":  payload["version"],
					"digest":   payload["digest"],
				},
			},
		}
	}
	_ = client
	return nil
}

func setGasData(transactionData *TransactionData, client CoreClient) error {
	if transactionData.GasData.Price == "" {
		var price string
		if err := client.Call(context.Background(), "suix_getReferenceGasPrice", []any{}, &price); err != nil {
			transactionData.GasData.Price = "1"
		} else {
			transactionData.GasData.Price = price
		}
	}
	if transactionData.GasData.Budget == "" {
		transactionData.GasData.Budget = itoa64(MaxGas)
	}
	if transactionData.GasData.Payment == nil {
		transactionData.GasData.Payment = []ObjectRef{}
	}
	if transactionData.Expiration == nil {
		transactionData.Expiration = map[string]any{
			"$kind": "ValidDuring",
			"ValidDuring": map[string]any{
				"minEpoch": "0",
				"maxEpoch": "1",
				"chain":    "unknown",
				"nonce":    0,
			},
		}
	}
	return nil
}

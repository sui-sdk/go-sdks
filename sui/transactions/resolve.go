package transactions

import (
	"context"
	"errors"
)

type BuildTransactionOptions struct {
	Client           CoreClient
	OnlyTransactionKind bool
}

type SerializeTransactionOptions struct {
	BuildTransactionOptions
	SupportedIntents []string
}

type TransactionPlugin func(transactionData *TransactionData, options BuildTransactionOptions, next func() error) error

type CoreClient interface {
	Call(ctx context.Context, method string, params []any, out any) error
}

func NeedsTransactionResolution(data *TransactionData, options BuildTransactionOptions) bool {
	for _, input := range data.Inputs {
		if input["$kind"] == "UnresolvedObject" || input["$kind"] == "UnresolvedPure" {
			return true
		}
	}
	if !options.OnlyTransactionKind {
		if data.GasData.Price == "" || data.GasData.Budget == "" || data.GasData.Payment == nil {
			return true
		}
		if len(data.GasData.Payment) == 0 && data.Expiration == nil {
			return true
		}
	}
	return false
}

func ResolveTransactionPlugin(transactionData *TransactionData, options BuildTransactionOptions, next func() error) error {
	normalizeRawArguments(transactionData)
	if !NeedsTransactionResolution(transactionData, options) {
		if err := validateResolvedInputs(transactionData); err != nil {
			return err
		}
		return next()
	}
	if options.Client == nil {
		return errors.New("no client passed to resolve transaction")
	}
	if err := CoreClientResolveTransaction(transactionData, options); err != nil {
		return err
	}
	if err := validateResolvedInputs(transactionData); err != nil {
		return err
	}
	return next()
}

func validateResolvedInputs(transactionData *TransactionData) error {
	for i, input := range transactionData.Inputs {
		kind, _ := input["$kind"].(string)
		if kind != "Object" && kind != "Pure" {
			return errors.New("input at index " + itoa(i) + " has not been resolved")
		}
	}
	return nil
}

func normalizeRawArguments(transactionData *TransactionData) {
	for _, cmd := range transactionData.Commands {
		kind, _ := cmd["$kind"].(string)
		switch kind {
		case "TransferObjects":
			payload, _ := cmd["TransferObjects"].(map[string]any)
			if addrArg, ok := payload["address"].(map[string]any); ok {
				normalizeRawArgument(addrArg, transactionData)
			}
		case "SplitCoins":
			payload, _ := cmd["SplitCoins"].(map[string]any)
			if amounts, ok := payload["amounts"].([]Argument); ok {
				for _, a := range amounts {
					normalizeRawArgument(a, transactionData)
				}
			}
		}
	}
}

func normalizeRawArgument(arg Argument, transactionData *TransactionData) {
	if arg["$kind"] != "Input" {
		return
	}
	idx, ok := arg["Input"].(int)
	if !ok || idx < 0 || idx >= len(transactionData.Inputs) {
		return
	}
	input := transactionData.Inputs[idx]
	if input["$kind"] != "UnresolvedPure" {
		return
	}
	value := input["UnresolvedPure"]
	// baseline fallback: encode as text bytes.
	transactionData.Inputs[idx] = Inputs.Pure([]byte(toString(value)))
}

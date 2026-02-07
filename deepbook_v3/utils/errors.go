package utils

import "fmt"

type DeepBookError struct{ Msg string }

func (e *DeepBookError) Error() string { return e.Msg }

type ResourceNotFoundError struct{ DeepBookError }
type ConfigurationError struct{ DeepBookError }
type ValidationError struct{ DeepBookError }

var ErrorMessages = struct {
	AdminCapNotSet           string
	MarginAdminCapNotSet     string
	MarginMaintainerCapNotSet string
	CoinNotFound             func(string) string
	PoolNotFound             func(string) string
	MarginPoolNotFound       func(string) string
	BalanceManagerNotFound   func(string) string
	MarginManagerNotFound    func(string) string
	PriceInfoNotFound        func(string) string
	InvalidAddress           string
}{
	AdminCapNotSet:            "Admin capability not configured",
	MarginAdminCapNotSet:      "Margin admin capability not configured",
	MarginMaintainerCapNotSet: "Margin maintainer capability not configured",
	CoinNotFound:              func(key string) string { return fmt.Sprintf("Coin not found for key: %s", key) },
	PoolNotFound:              func(key string) string { return fmt.Sprintf("Pool not found for key: %s", key) },
	MarginPoolNotFound:        func(key string) string { return fmt.Sprintf("Margin pool not found for key: %s", key) },
	BalanceManagerNotFound:    func(key string) string { return fmt.Sprintf("Balance manager with key %s not found", key) },
	MarginManagerNotFound:     func(key string) string { return fmt.Sprintf("Margin manager with key %s not found", key) },
	PriceInfoNotFound:         func(key string) string { return fmt.Sprintf("Price info object not found for %s", key) },
	InvalidAddress:            "Address must be a valid Sui address",
}

package jsonrpc

import "fmt"

func GetJSONRPCFullnodeURL(network string) (string, error) {
	switch network {
	case "mainnet":
		return "https://fullnode.mainnet.sui.io:443", nil
	case "testnet":
		return "https://fullnode.testnet.sui.io:443", nil
	case "devnet":
		return "https://fullnode.devnet.sui.io:443", nil
	case "localnet":
		return "http://127.0.0.1:9000", nil
	default:
		return "", fmt.Errorf("unknown network: %s", network)
	}
}

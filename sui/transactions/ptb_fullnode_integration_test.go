package transactions

import (
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/sui-sdks/go-sdks/sui/faucet"
	jsonrpc "github.com/sui-sdks/go-sdks/sui/jsonrpc"
	ed25519 "github.com/sui-sdks/go-sdks/sui/keypairs/ed25519"
)

type fullnodeITConfig struct {
	network       string
	fullnodeURL   string
	enableExecute bool
	privateKey    string
	faucetURL     string
}

func loadFullnodeITConfig(t *testing.T) fullnodeITConfig {
	t.Helper()
	if os.Getenv("SUI_IT_ENABLE_FULLNODE") != "1" {
		t.Skip("set SUI_IT_ENABLE_FULLNODE=1 to run real fullnode PTB integration tests")
	}

	network := getenv("SUI_IT_NETWORK", "testnet")
	fullnodeURL := os.Getenv("SUI_IT_FULLNODE_URL")
	if fullnodeURL == "" {
		u, err := jsonrpc.GetJSONRPCFullnodeURL(network)
		if err != nil {
			t.Fatalf("resolve fullnode url failed: %v", err)
		}
		fullnodeURL = u
	}

	faucetURL := os.Getenv("SUI_IT_FAUCET_URL")
	if faucetURL == "" {
		if u, err := faucet.GetFaucetHost(network); err == nil {
			faucetURL = u
		}
	}

	return fullnodeITConfig{
		network:       network,
		fullnodeURL:   fullnodeURL,
		enableExecute: os.Getenv("SUI_IT_ENABLE_EXECUTE") == "1",
		privateKey:    os.Getenv("SUI_IT_PRIVATE_KEY"),
		faucetURL:     faucetURL,
	}
}

func newITJSONRPCClient(t *testing.T, cfg fullnodeITConfig) *jsonrpc.Client {
	t.Helper()
	c, err := jsonrpc.NewClient(jsonrpc.ClientOptions{
		Network: cfg.network,
		URL:     cfg.fullnodeURL,
	})
	if err != nil {
		t.Fatalf("new jsonrpc client failed: %v", err)
	}
	return c
}

func TestPTBFullnodeIntegrationResolveBuild(t *testing.T) {
	cfg := loadFullnodeITConfig(t)
	client := newITJSONRPCClient(t, cfg)

	var gasPrice string
	if err := client.Call(context.Background(), "suix_getReferenceGasPrice", []any{}, &gasPrice); err != nil {
		t.Fatalf("fullnode health check failed: %v", err)
	}
	if gasPrice == "" {
		t.Fatalf("unexpected empty gas price")
	}

	tx := NewTransaction()
	tx.SetSender(getenv("SUI_IT_SENDER", "0x1"))

	split := tx.SplitCoins(tx.Gas(), []Argument{tx.PureBytes([]byte{1})})
	tx.TransferObjects([]Argument{split}, tx.PureBytes([]byte("0x2")))

	exec := NewCachingTransactionExecutor(client, nil)
	built, err := exec.BuildTransaction(tx, BuildTransactionOptions{Client: client})
	if err != nil {
		t.Fatalf("build transaction with fullnode resolver failed: %v", err)
	}
	if len(built) == 0 {
		t.Fatalf("expected built transaction bytes")
	}

	restored, err := TransactionFrom(base64.StdEncoding.EncodeToString(built))
	if err != nil {
		t.Fatalf("restore transaction from built bytes failed: %v", err)
	}
	data := restored.GetData()
	if data.GasData.Price == "" || data.GasData.Budget == "" {
		t.Fatalf("expected gas data to be resolved, got price=%q budget=%q", data.GasData.Price, data.GasData.Budget)
	}
	if data.Expiration == nil {
		t.Fatalf("expected expiration to be resolved")
	}
	if len(data.Commands) < 2 {
		t.Fatalf("expected at least 2 commands, got %d", len(data.Commands))
	}
}

func TestPTBFullnodeIntegrationSignAndExecute(t *testing.T) {
	cfg := loadFullnodeITConfig(t)
	if !cfg.enableExecute {
		t.Skip("set SUI_IT_ENABLE_EXECUTE=1 to run sign+execute integration test")
	}
	if cfg.privateKey == "" {
		t.Skip("set SUI_IT_PRIVATE_KEY=suiprivkey:... to run sign+execute integration test")
	}

	client := newITJSONRPCClient(t, cfg)
	signer, err := ed25519.FromSecretKeyString(cfg.privateKey)
	if err != nil {
		t.Fatalf("parse SUI_IT_PRIVATE_KEY failed: %v", err)
	}
	ensureFunds(t, client, cfg, signer.ToSuiAddress())

	tx := NewTransaction()
	tx.SetSender(signer.ToSuiAddress())
	tx.SetGasBudget(50_000_000)

	split := tx.SplitCoins(tx.Gas(), []Argument{tx.PureBytes([]byte{1})})
	tx.TransferObjects([]Argument{split}, tx.PureBytes([]byte(signer.ToSuiAddress())))

	exec := NewSerialTransactionExecutor(SerialTransactionExecutorOptions{
		Client: client,
		Signer: signer,
	})
	res, err := exec.ExecuteTransaction(tx, map[string]any{
		"showEffects":        true,
		"showBalanceChanges": true,
	}, nil)
	if err != nil {
		t.Fatalf("execute PTB transaction failed: %v", err)
	}

	digest, _ := res["digest"].(string)
	if digest == "" {
		t.Fatalf("missing digest in execute response: %+v", res)
	}

	if _, err := client.GetTransactionBlock(context.Background(), digest, map[string]any{
		"showEffects": true,
	}); err != nil {
		t.Fatalf("query executed transaction failed: %v", err)
	}
}

func ensureFunds(t *testing.T, client *jsonrpc.Client, cfg fullnodeITConfig, address string) {
	t.Helper()
	balance, err := client.GetBalance(context.Background(), address, "")
	if err == nil {
		if tb, ok := balance["totalBalance"].(string); ok && tb != "" && tb != "0" {
			return
		}
	}

	if cfg.faucetURL == "" {
		t.Fatalf("address %s has no funds and no faucet configured (set SUI_IT_FAUCET_URL)", address)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	_, err = faucet.RequestSuiFromFaucetV2(ctx, cfg.faucetURL, address, nil)
	if err != nil {
		t.Fatalf("request faucet funds failed: %v", err)
	}

	deadline := time.Now().Add(45 * time.Second)
	for time.Now().Before(deadline) {
		time.Sleep(3 * time.Second)
		b, err := client.GetBalance(context.Background(), address, "")
		if err != nil {
			continue
		}
		if tb, ok := b["totalBalance"].(string); ok && tb != "" && tb != "0" {
			return
		}
	}
	t.Fatalf("faucet funds not visible in time for address %s", address)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}


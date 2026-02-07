package pyth

import (
	"context"
	"encoding/base64"
	"errors"

	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

type PriceServiceConnection interface {
	GetLatestVaas(priceIDs []string) ([]string, error)
}

type SuiPriceServiceConnection struct {
	inner PriceServiceConnection
}

func NewSuiPriceServiceConnection(inner PriceServiceConnection) *SuiPriceServiceConnection {
	return &SuiPriceServiceConnection{inner: inner}
}

func (c *SuiPriceServiceConnection) GetPriceFeedsUpdateData(priceIDs []string) ([][]byte, error) {
	latest, err := c.inner.GetLatestVaas(priceIDs)
	if err != nil {
		return nil, err
	}
	out := make([][]byte, 0, len(latest))
	for _, v := range latest {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

type CoreClient interface {
	Call(ctx context.Context, method string, params []any, out any) error
}

type SuiPythClient struct {
	provider        CoreClient
	PythStateID     string
	WormholeStateID string
}

func NewSuiPythClient(provider CoreClient, pythStateID, wormholeStateID string) *SuiPythClient {
	return &SuiPythClient{provider: provider, PythStateID: pythStateID, WormholeStateID: wormholeStateID}
}

func (c *SuiPythClient) VerifyVaas(tx *stx.Transaction, vaas [][]byte) []stx.Argument {
	out := make([]stx.Argument, 0, len(vaas))
	for _, vaa := range vaas {
		arg := tx.MoveCall("wormhole::vaa::parse_and_verify", []stx.Argument{tx.Object(c.WormholeStateID), tx.PureBytes(vaa), tx.Object("0x6")}, nil)
		out = append(out, arg)
	}
	return out
}

func (c *SuiPythClient) UpdatePriceFeeds(tx *stx.Transaction, updates [][]byte, feedIDs []string) ([]string, error) {
	if len(updates) > 1 {
		return nil, errors.New("multiple accumulator messages are not supported")
	}
	_ = tx
	return append([]string(nil), feedIDs...), nil
}

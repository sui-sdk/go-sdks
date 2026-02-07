package transactions

import "testing"

func TestTransactionBuildSerialize(t *testing.T) {
	tx := NewTransaction()
	tx.SetSender("0x1")
	tx.SetGasBudget(1000)
	tx.SetGasPrice(1)
	coin := tx.Gas()
	a1 := tx.PureBytes([]byte{10})
	res := tx.SplitCoins(coin, []Argument{a1})
	addr := tx.PureBytes([]byte("addr"))
	tx.TransferObjects([]Argument{res}, addr)

	b, err := tx.Build()
	if err != nil {
		t.Fatalf("build failed: %v", err)
	}
	if len(b) == 0 {
		t.Fatalf("expected bytes")
	}
	serialized, err := tx.Serialize()
	if err != nil {
		t.Fatalf("serialize failed: %v", err)
	}
	restored, err := TransactionFrom(serialized)
	if err != nil {
		t.Fatalf("restore failed: %v", err)
	}
	if restored.GetData().Sender == "" {
		t.Fatalf("missing sender")
	}
}

package transactions

import (
	"encoding/base64"
	"testing"
)

func TestPTBBasicsSplitAndTransfer(t *testing.T) {
	tx := NewTransaction()
	tx.SetSender("0x1")
	tx.SetGasBudget(1_000_000)
	tx.SetGasPrice(1)

	coin := tx.SplitCoins(tx.Gas(), []Argument{
		tx.PureBytes([]byte{100}),
	})
	recipient := tx.PureBytes([]byte("0x2"))
	tx.TransferObjects([]Argument{coin}, recipient)
	tx.MergeCoins(tx.Gas(), []Argument{coin})

	data := tx.GetData()
	if len(data.Commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(data.Commands))
	}
	if got := data.Commands[0]["$kind"]; got != "SplitCoins" {
		t.Fatalf("expected first command SplitCoins, got %v", got)
	}
	if got := data.Commands[1]["$kind"]; got != "TransferObjects" {
		t.Fatalf("expected second command TransferObjects, got %v", got)
	}
	if got := data.Commands[2]["$kind"]; got != "MergeCoins" {
		t.Fatalf("expected third command MergeCoins, got %v", got)
	}
}

func TestPTBBasicsMoveCallAndMakeMoveVec(t *testing.T) {
	tx := NewTransaction()
	obj := tx.Object("0x2")
	amount := tx.PureBytes([]byte{1, 2, 3})

	tx.MoveCall(
		"0x2::example::do_something",
		[]Argument{obj, amount},
		[]string{"0x2::sui::SUI"},
	)

	vec := tx.AddCommand(TransactionCommands.MakeMoveVec(nil, []Argument{
		tx.Object("0x3"),
		tx.Object("0x4"),
	}))
	tx.TransferObjects([]Argument{vec}, tx.PureBytes([]byte("0x5")))

	data := tx.GetData()
	if len(data.Commands) != 3 {
		t.Fatalf("expected 3 commands, got %d", len(data.Commands))
	}

	moveCall := data.Commands[0]["MoveCall"].(map[string]any)
	if moveCall["package"] != "0x2" || moveCall["module"] != "example" || moveCall["function"] != "do_something" {
		t.Fatalf("unexpected move call target parsing: %+v", moveCall)
	}

	if got := data.Commands[1]["$kind"]; got != "MakeMoveVec" {
		t.Fatalf("expected second command MakeMoveVec, got %v", got)
	}
}

func TestPTBBasicsPublishAndResolvedObjectRefs(t *testing.T) {
	tx := NewTransaction()
	tx.Publish(
		[][]byte{[]byte{0xAA, 0xBB}},
		[]string{"0x2", "0x0003"},
	)
	tx.Object(Inputs.ObjectRef(ObjectRef{
		ObjectID: "0x123",
		Digest:   "digest1",
		Version:  "7",
	}))
	tx.Object(Inputs.SharedObjectRef("0x124", true, "10"))
	tx.Object(Inputs.ReceivingRef(ObjectRef{
		ObjectID: "0x125",
		Digest:   "digest2",
		Version:  "11",
	}))

	data := tx.GetData()
	if len(data.Commands) != 1 {
		t.Fatalf("expected 1 command, got %d", len(data.Commands))
	}
	pub := data.Commands[0]["Publish"].(map[string]any)
	mods := pub["modules"].([]string)
	if len(mods) != 1 || mods[0] != base64.StdEncoding.EncodeToString([]byte{0xAA, 0xBB}) {
		t.Fatalf("unexpected publish modules: %+v", mods)
	}
	if len(data.Inputs) != 3 {
		t.Fatalf("expected 3 resolved object inputs, got %d", len(data.Inputs))
	}
}

func TestPTBBasicsBuildBase64AndFrom(t *testing.T) {
	tx := NewTransaction()
	tx.SetSender("0x1")
	tx.SplitCoins(tx.Gas(), []Argument{tx.PureBytes([]byte{9})})

	b64, err := tx.BuildBase64()
	if err != nil {
		t.Fatalf("build base64 failed: %v", err)
	}
	restored, err := TransactionFrom(b64)
	if err != nil {
		t.Fatalf("transaction from base64 failed: %v", err)
	}
	if len(restored.GetData().Commands) != 1 {
		t.Fatalf("expected restored command count 1, got %d", len(restored.GetData().Commands))
	}
}

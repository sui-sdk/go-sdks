package transactions

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/sui-sdks/go-sdks/sui/utils"
)

type GasData struct {
	Owner   string
	Price   string
	Budget  string
	Payment []ObjectRef
}

type TransactionData struct {
	Sender     string
	Expiration any
	GasData    GasData
	Inputs     []CallArg
	Commands   []Command
}

type Transaction struct {
	data TransactionData
}

func NewTransaction() *Transaction {
	return &Transaction{data: TransactionData{Inputs: []CallArg{}, Commands: []Command{}, GasData: GasData{}}}
}

func TransactionFrom(serialized string) (*Transaction, error) {
	if len(serialized) > 0 && serialized[0] == '{' {
		var data TransactionData
		if err := json.Unmarshal([]byte(serialized), &data); err != nil {
			return nil, err
		}
		return &Transaction{data: data}, nil
	}
	b, err := base64.StdEncoding.DecodeString(serialized)
	if err != nil {
		return nil, err
	}
	var data TransactionData
	if err := json.Unmarshal(b, &data); err != nil {
		return nil, err
	}
	return &Transaction{data: data}, nil
}

func (t *Transaction) GetData() TransactionData { return t.data }
func (t *Transaction) SetSender(sender string) { t.data.Sender = utils.NormalizeSuiAddress(sender) }
func (t *Transaction) SetSenderIfNotSet(sender string) {
	if t.data.Sender == "" {
		t.SetSender(sender)
	}
}
func (t *Transaction) SetExpiration(exp any) { t.data.Expiration = exp }
func (t *Transaction) SetGasPrice(price any) { t.data.GasData.Price = fmt.Sprintf("%v", price) }
func (t *Transaction) SetGasBudget(budget any) { t.data.GasData.Budget = fmt.Sprintf("%v", budget) }
func (t *Transaction) SetGasBudgetIfNotSet(budget any) {
	if t.data.GasData.Budget == "" {
		t.SetGasBudget(budget)
	}
}
func (t *Transaction) SetGasOwner(owner string) { t.data.GasData.Owner = utils.NormalizeSuiAddress(owner) }
func (t *Transaction) SetGasPayment(payments []ObjectRef) { t.data.GasData.Payment = payments }

func (t *Transaction) Gas() Argument { return Argument{"$kind": "GasCoin", "GasCoin": true} }

func (t *Transaction) AddInput(arg CallArg) Argument {
	t.data.Inputs = append(t.data.Inputs, arg)
	return Argument{"$kind": "Input", "Input": len(t.data.Inputs) - 1}
}

func (t *Transaction) Object(value any) Argument {
	switch v := value.(type) {
	case string:
		return t.AddInput(CallArg{"$kind": "UnresolvedObject", "UnresolvedObject": map[string]any{"objectId": utils.NormalizeSuiObjectID(v)}})
	case CallArg:
		return t.AddInput(v)
	case Argument:
		return v
	default:
		return t.AddInput(CallArg{"$kind": "UnresolvedObject", "UnresolvedObject": map[string]any{"value": v}})
	}
}

func (t *Transaction) PureBytes(value []byte) Argument {
	return t.AddInput(Inputs.Pure(value))
}

func (t *Transaction) AddCommand(cmd Command) Argument {
	t.data.Commands = append(t.data.Commands, cmd)
	return Argument{"$kind": "Result", "Result": len(t.data.Commands) - 1}
}

func (t *Transaction) MoveCall(target string, args []Argument, typeArgs []string) Argument {
	return t.AddCommand(TransactionCommands.MoveCall(target, args, typeArgs))
}

func (t *Transaction) TransferObjects(objects []Argument, address Argument) Argument {
	return t.AddCommand(TransactionCommands.TransferObjects(objects, address))
}

func (t *Transaction) SplitCoins(coin Argument, amounts []Argument) Argument {
	return t.AddCommand(TransactionCommands.SplitCoins(coin, amounts))
}

func (t *Transaction) MergeCoins(destination Argument, sources []Argument) Argument {
	return t.AddCommand(TransactionCommands.MergeCoins(destination, sources))
}

func (t *Transaction) Publish(modules [][]byte, dependencies []string) Argument {
	return t.AddCommand(TransactionCommands.Publish(modules, dependencies))
}

func (t *Transaction) Build() ([]byte, error) {
	return json.Marshal(t.data)
}

func (t *Transaction) BuildBase64() (string, error) {
	b, err := t.Build()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func (t *Transaction) Serialize() (string, error) {
	b, err := json.Marshal(t.data)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

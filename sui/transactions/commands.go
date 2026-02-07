package transactions

import (
	"encoding/base64"
	"strings"

	"github.com/sui-sdks/go-sdks/sui/utils"
)

type Argument map[string]any

type Command map[string]any

const (
	UpgradePolicyCompatible = 0
	UpgradePolicyAdditive   = 128
	UpgradePolicyDepOnly    = 192
)

var TransactionCommands = struct {
	MoveCall        func(target string, args []Argument, typeArgs []string) Command
	TransferObjects func(objects []Argument, address Argument) Command
	SplitCoins      func(coin Argument, amounts []Argument) Command
	MergeCoins      func(destination Argument, sources []Argument) Command
	Publish         func(modules [][]byte, dependencies []string) Command
	MakeMoveVec     func(typ *string, elements []Argument) Command
}{
	MoveCall: func(target string, args []Argument, typeArgs []string) Command {
		parts := strings.Split(target, "::")
		pkg, mod, fn := "", "", ""
		if len(parts) > 0 {
			pkg = parts[0]
		}
		if len(parts) > 1 {
			mod = parts[1]
		}
		if len(parts) > 2 {
			fn = parts[2]
		}
		return Command{"$kind": "MoveCall", "MoveCall": map[string]any{"package": pkg, "module": mod, "function": fn, "arguments": args, "typeArguments": typeArgs}}
	},
	TransferObjects: func(objects []Argument, address Argument) Command {
		return Command{"$kind": "TransferObjects", "TransferObjects": map[string]any{"objects": objects, "address": address}}
	},
	SplitCoins: func(coin Argument, amounts []Argument) Command {
		return Command{"$kind": "SplitCoins", "SplitCoins": map[string]any{"coin": coin, "amounts": amounts}}
	},
	MergeCoins: func(destination Argument, sources []Argument) Command {
		return Command{"$kind": "MergeCoins", "MergeCoins": map[string]any{"destination": destination, "sources": sources}}
	},
	Publish: func(modules [][]byte, dependencies []string) Command {
		mods := make([]string, len(modules))
		for i := range modules {
			mods[i] = base64.StdEncoding.EncodeToString(modules[i])
		}
		deps := make([]string, len(dependencies))
		for i := range dependencies {
			deps[i] = utils.NormalizeSuiObjectID(dependencies[i])
		}
		return Command{"$kind": "Publish", "Publish": map[string]any{"modules": mods, "dependencies": deps}}
	},
	MakeMoveVec: func(typ *string, elements []Argument) Command {
		var t any
		if typ != nil {
			t = *typ
		}
		return Command{"$kind": "MakeMoveVec", "MakeMoveVec": map[string]any{"type": t, "elements": elements}}
	},
}

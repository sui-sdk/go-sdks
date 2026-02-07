package transactions

var Arguments = struct {
	Pure      func([]byte) func(*Transaction) Argument
	Object    func(any) func(*Transaction) Argument
	ObjectRef func(ObjectRef) func(*Transaction) Argument
}{
	Pure: func(v []byte) func(*Transaction) Argument {
		return func(tx *Transaction) Argument { return tx.PureBytes(v) }
	},
	Object: func(v any) func(*Transaction) Argument {
		return func(tx *Transaction) Argument { return tx.Object(v) }
	},
	ObjectRef: func(ref ObjectRef) func(*Transaction) Argument {
		return func(tx *Transaction) Argument { return tx.Object(Inputs.ObjectRef(ref)) }
	},
}

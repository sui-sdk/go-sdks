package transactions

import (
	"encoding/base64"

	"github.com/sui-sdks/go-sdks/sui/utils"
)

type ObjectRef struct {
	ObjectID string `json:"objectId"`
	Digest   string `json:"digest"`
	Version  any    `json:"version"`
}

type CallArg map[string]any

var Inputs = struct {
	Pure            func([]byte) CallArg
	ObjectRef       func(ObjectRef) CallArg
	SharedObjectRef func(objectID string, mutable bool, initialSharedVersion any) CallArg
	ReceivingRef    func(ObjectRef) CallArg
}{
	Pure: func(data []byte) CallArg {
		return CallArg{"$kind": "Pure", "Pure": map[string]any{"bytes": base64.StdEncoding.EncodeToString(data)}}
	},
	ObjectRef: func(ref ObjectRef) CallArg {
		return CallArg{"$kind": "Object", "Object": map[string]any{"$kind": "ImmOrOwnedObject", "ImmOrOwnedObject": map[string]any{"objectId": utils.NormalizeSuiAddress(ref.ObjectID), "digest": ref.Digest, "version": ref.Version}}}
	},
	SharedObjectRef: func(objectID string, mutable bool, initialSharedVersion any) CallArg {
		return CallArg{"$kind": "Object", "Object": map[string]any{"$kind": "SharedObject", "SharedObject": map[string]any{"objectId": utils.NormalizeSuiAddress(objectID), "mutable": mutable, "initialSharedVersion": initialSharedVersion}}}
	},
	ReceivingRef: func(ref ObjectRef) CallArg {
		return CallArg{"$kind": "Object", "Object": map[string]any{"$kind": "Receiving", "Receiving": map[string]any{"objectId": utils.NormalizeSuiAddress(ref.ObjectID), "digest": ref.Digest, "version": ref.Version}}}
	},
}

package zklogin

import (
	"encoding/base64"
	"encoding/json"
)

type ZkLoginSignatureExtended struct {
	Inputs       map[string]any `json:"inputs"`
	MaxEpoch     uint64         `json:"maxEpoch"`
	UserSignature string        `json:"userSignature"`
}

func GetZkLoginSignature(input ZkLoginSignatureExtended) (string, error) {
	b, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func ParseZkLoginSignature(signature string) (ZkLoginSignatureExtended, error) {
	b, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return ZkLoginSignatureExtended{}, err
	}
	var out ZkLoginSignatureExtended
	err = json.Unmarshal(b, &out)
	return out, err
}

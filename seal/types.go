package seal

import "time"

type KeyCacheKey string

type SealOptions struct {
	Name            string
	ServerConfigs   []KeyServerConfig
	VerifyKeyServers *bool
	Timeout         time.Duration
}

type KeyServerConfig struct {
	ObjectID      string
	Weight        int
	APIKeyName    string
	APIKey        string
	AggregatorURL string
	URL           string
	PublicKeyHex  string
}

type ClientOptions struct {
	ServerConfigs   []KeyServerConfig
	VerifyKeyServers *bool
	Timeout         time.Duration
}

type EncryptOptions struct {
	KEMType   KemType
	DEMType   DemType
	Threshold int
	PackageID string
	ID        string
	Data      []byte
	AAD       []byte
}

type DecryptOptions struct {
	Data                   []byte
	SessionKey             *SessionKey
	TxBytes                []byte
	CheckShareConsistency  bool
	CheckLEEncoding        bool
}

type FetchKeysOptions struct {
	IDs       []string
	TxBytes   []byte
	SessionKey *SessionKey
	Threshold int
}

type GetDerivedKeysOptions struct {
	KEMType   KemType
	ID        string
	TxBytes   []byte
	SessionKey *SessionKey
	Threshold int
}

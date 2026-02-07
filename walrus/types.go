package walrus

type TipStrategyConst struct {
	Const uint64
}

type TipStrategyLinear struct {
	Base         uint64
	PerEncodedKiB uint64
}

type UploadRelayTipConfig struct {
	Address string
	Max     *int
	Kind    any
}

type UploadRelayConfig struct {
	URL     string
	SendTip any
}

type Config struct {
	Network                 string
	PackageConfig           *PackageConfig
	StorageNodeClientConfig *StorageNodeClientOptions
	WasmURL                 string
	UploadRelay             *UploadRelayConfig
}

type ReadBlobOptions struct {
	BlobID string
}

type GetBlobMetadataOptions struct {
	BlobID string
}

type GetBlobStatusOptions struct {
	BlobID string
}

type GetSliverOptions struct {
	BlobID         string
	SliverPairIndex int
	SliverType     string
}

type WriteBlobMetadataOptions struct {
	BlobID   string
	Metadata []byte
	NodeURL  string
}

type WriteSliverOptions struct {
	BlobID         string
	SliverPairIndex int
	SliverType     string
	Sliver         []byte
	NodeURL        string
}

// Mirror storage-node options at walrus package level for convenience.
type StorageNodeClientOptions = storageNodeClientOptionsAlias

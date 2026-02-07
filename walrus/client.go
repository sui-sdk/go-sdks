package walrus

import (
	"context"
	"fmt"
)

type Client struct {
	storageNodeClient *StorageNodeClient
	packageConfig     PackageConfig
	network           string
}

type ExtensionOptions struct {
	Name          string
	PackageConfig *PackageConfig
}

func Walrus(opts *ExtensionOptions) *ExtensionOptions {
	if opts == nil {
		opts = &ExtensionOptions{}
	}
	if opts.Name == "" {
		opts.Name = "walrus"
	}
	return opts
}

func NewClient(config Config) (*Client, error) {
	pkg := config.PackageConfig
	if pkg == nil {
		switch config.Network {
		case "testnet":
			pc := TestnetWalrusPackageConfig
			pkg = &pc
		case "mainnet":
			pc := MainnetWalrusPackageConfig
			pkg = &pc
		default:
			return nil, fmt.Errorf("unsupported network: %s", config.Network)
		}
	}
	return &Client{
		storageNodeClient: NewStorageNodeClient((*storageNodeClientOptionsAlias)(config.StorageNodeClientConfig)),
		packageConfig:     *pkg,
		network:           config.Network,
	}, nil
}

func (c *Client) PackageConfig() PackageConfig {
	return c.packageConfig
}

func (c *Client) GetBlobMetadata(ctx context.Context, nodeURL string, opts GetBlobMetadataOptions) ([]byte, error) {
	return c.storageNodeClient.GetBlobMetadata(map[string]string{"blobId": opts.BlobID}, RequestOptions{Ctx: ctx, NodeURL: nodeURL})
}

func (c *Client) GetBlobStatus(ctx context.Context, nodeURL string, opts GetBlobStatusOptions) (map[string]any, error) {
	return c.storageNodeClient.GetBlobStatus(map[string]string{"blobId": opts.BlobID}, RequestOptions{Ctx: ctx, NodeURL: nodeURL})
}

func (c *Client) GetSliver(ctx context.Context, nodeURL string, opts GetSliverOptions) ([]byte, error) {
	return c.storageNodeClient.GetSliver(map[string]any{
		"blobId":         opts.BlobID,
		"sliverPairIndex": opts.SliverPairIndex,
		"sliverType":     opts.SliverType,
	}, RequestOptions{Ctx: ctx, NodeURL: nodeURL})
}

func (c *Client) ReadBlob(ctx context.Context, nodeURL string, opts ReadBlobOptions) ([]byte, error) {
	status, err := c.GetBlobStatus(ctx, nodeURL, GetBlobStatusOptions{BlobID: opts.BlobID})
	if err != nil {
		return nil, err
	}
	if status == nil {
		return nil, &NoBlobStatusReceivedError{ClientError{Msg: "no blob status received"}}
	}
	if _, ok := status["invalid"]; ok {
		return nil, &BlobBlockedError{Msg: "blob status invalid"}
	}
	// This method returns raw metadata bytes as a baseline implementation.
	return c.GetBlobMetadata(ctx, nodeURL, GetBlobMetadataOptions{BlobID: opts.BlobID})
}

func (c *Client) WriteBlobMetadata(ctx context.Context, opts WriteBlobMetadataOptions) (map[string]any, error) {
	return c.storageNodeClient.StoreBlobMetadata(map[string]any{
		"blobId":   opts.BlobID,
		"metadata": opts.Metadata,
	}, RequestOptions{Ctx: ctx, NodeURL: opts.NodeURL})
}

func (c *Client) WriteSliver(ctx context.Context, opts WriteSliverOptions) (map[string]any, error) {
	return c.storageNodeClient.StoreSliver(map[string]any{
		"blobId":         opts.BlobID,
		"sliverPairIndex": opts.SliverPairIndex,
		"sliverType":     opts.SliverType,
		"sliver":         opts.Sliver,
	}, RequestOptions{Ctx: ctx, NodeURL: opts.NodeURL})
}

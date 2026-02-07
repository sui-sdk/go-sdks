package seal

import (
	"context"
	"fmt"
	"time"
)

type Extension struct {
	Name string
	Opts SealOptions
}

func Seal(opts SealOptions) Extension {
	if opts.Name == "" {
		opts.Name = "seal"
	}
	return Extension{Name: opts.Name, Opts: opts}
}

type Client struct {
	configs          map[string]KeyServerConfig
	keyServers       map[string]KeyServer
	verifyKeyServers bool
	cachedKeys       map[KeyCacheKey][]byte
	timeout          time.Duration
	totalWeight      int
}

func NewClient(options ClientOptions) (*Client, error) {
	if len(options.ServerConfigs) == 0 {
		return nil, (&InvalidClientOptionsError{UserError{SealError{Msg: "serverConfigs is required"}}})
	}
	verify := true
	if options.VerifyKeyServers != nil {
		verify = *options.VerifyKeyServers
	}
	timeout := options.Timeout
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	cfgs := make(map[string]KeyServerConfig, len(options.ServerConfigs))
	totalWeight := 0
	for _, s := range options.ServerConfigs {
		if _, ok := cfgs[s.ObjectID]; ok {
			return nil, (&InvalidClientOptionsError{UserError{SealError{Msg: "duplicate object IDs"}}})
		}
		if (s.APIKeyName == "") != (s.APIKey == "") {
			return nil, (&InvalidClientOptionsError{UserError{SealError{Msg: "apiKeyName and apiKey must be set together"}}})
		}
		cfgs[s.ObjectID] = s
		totalWeight += s.Weight
	}
	return &Client{configs: cfgs, verifyKeyServers: verify, cachedKeys: map[KeyCacheKey][]byte{}, timeout: timeout, totalWeight: totalWeight}, nil
}

func (c *Client) Encrypt(opts EncryptOptions) (encryptedObject []byte, key []byte, err error) {
	if opts.KEMType != KemTypeBonehFranklinBLS12381DemCCA {
		opts.KEMType = KemTypeBonehFranklinBLS12381DemCCA
	}
	if opts.DEMType != DemTypeHmac256Ctr {
		opts.DEMType = DemTypeAesGcm256
	}
	servers, err := c.getWeightedKeyServers()
	if err != nil {
		return nil, nil, err
	}
	var encInput EncryptionInput
	switch opts.DEMType {
	case DemTypeAesGcm256:
		encInput = AesGcm256{Data: opts.Data, AAD: opts.AAD}
	case DemTypeHmac256Ctr:
		encInput = Hmac256Ctr{Data: opts.Data, AAD: opts.AAD}
	}
	return Encrypt(struct {
		KeyServers      []KeyServer
		KEMType         KemType
		Threshold       int
		PackageID       string
		ID              string
		EncryptionInput EncryptionInput
	}{KeyServers: servers, KEMType: opts.KEMType, Threshold: opts.Threshold, PackageID: opts.PackageID, ID: opts.ID, EncryptionInput: encInput})
}

func (c *Client) Decrypt(opts DecryptOptions) ([]byte, error) {
	obj, err := ParseEncryptedObject(opts.Data)
	if err != nil {
		return nil, err
	}
	if err := c.FetchKeys(context.Background(), FetchKeysOptions{IDs: []string{obj.ID}, TxBytes: opts.TxBytes, SessionKey: opts.SessionKey, Threshold: obj.Threshold}); err != nil {
		return nil, err
	}
	return Decrypt(struct {
		EncryptedObject EncryptedObject
		Keys            map[KeyCacheKey][]byte
		CheckLEEncoding bool
	}{EncryptedObject: obj, Keys: c.cachedKeys, CheckLEEncoding: opts.CheckLEEncoding})
}

func (c *Client) GetKeyServers() (map[string]KeyServer, error) {
	if c.keyServers == nil {
		servers, err := RetrieveKeyServers(c.configs)
		if err != nil {
			return nil, err
		}
		m := make(map[string]KeyServer, len(servers))
		for _, s := range servers {
			if c.verifyKeyServers {
				cfg := c.configs[s.ObjectID]
				if !VerifyKeyServer(s, c.timeout, cfg.APIKeyName, cfg.APIKey) {
					return nil, (&InvalidKeyServerError{UserError{SealError{Msg: fmt.Sprintf("key server %s is not valid", s.ObjectID)}}})
				}
			}
			m[s.ObjectID] = s
		}
		c.keyServers = m
	}
	return c.keyServers, nil
}

func (c *Client) FetchKeys(ctx context.Context, opts FetchKeysOptions) error {
	if opts.Threshold > c.totalWeight || opts.Threshold < 1 {
		return (&InvalidThresholdError{UserError{SealError{Msg: fmt.Sprintf("invalid threshold %d", opts.Threshold)}}})
	}
	servers, err := c.GetKeyServers()
	if err != nil {
		return err
	}
	fullIDs := make([]string, len(opts.IDs))
	for i, id := range opts.IDs {
		fullIDs[i] = CreateFullID(opts.SessionKey.GetPackageID(), id)
	}

	completedWeight := 0
	for objectID, cfg := range c.configs {
		complete := true
		for _, fullID := range fullIDs {
			if _, ok := c.cachedKeys[KeyCacheKey(fullID+":"+objectID)]; !ok {
				complete = false
				break
			}
		}
		if complete {
			completedWeight += cfg.Weight
		}
	}
	if completedWeight >= opts.Threshold {
		return nil
	}

	errors := make([]error, 0)
	for objectID, server := range servers {
		cfg := c.configs[objectID]
		derived, err := FetchKeysForAllIDs(ctx, server, fullIDs, opts.TxBytes, opts.SessionKey, c.timeout, cfg.APIKeyName, cfg.APIKey)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		for _, d := range derived {
			c.cachedKeys[KeyCacheKey(d.ID+":"+d.ObjectID)] = d.Key
		}
		completedWeight += cfg.Weight
		if completedWeight >= opts.Threshold {
			return nil
		}
	}
	if completedWeight < opts.Threshold {
		if len(errors) > 0 {
			return (&TooManyFailedFetchKeyRequestsError{UserError{SealError{Msg: ToMajorityError(errors).Error()}}})
		}
		return (&TooManyFailedFetchKeyRequestsError{UserError{SealError{Msg: "threshold not reached"}}})
	}
	return nil
}

func (c *Client) GetDerivedKeys(ctx context.Context, opts GetDerivedKeysOptions) (map[string][]byte, error) {
	if err := c.FetchKeys(ctx, FetchKeysOptions{IDs: []string{opts.ID}, TxBytes: opts.TxBytes, SessionKey: opts.SessionKey, Threshold: opts.Threshold}); err != nil {
		return nil, err
	}
	fullID := CreateFullID(opts.SessionKey.GetPackageID(), opts.ID)
	result := map[string][]byte{}
	for objectID := range c.configs {
		if key, ok := c.cachedKeys[KeyCacheKey(fullID+":"+objectID)]; ok {
			result[objectID] = append([]byte(nil), key...)
		}
	}
	return result, nil
}

func (c *Client) getWeightedKeyServers() ([]KeyServer, error) {
	servers, err := c.GetKeyServers()
	if err != nil {
		return nil, err
	}
	out := make([]KeyServer, 0)
	for id, cfg := range c.configs {
		s := servers[id]
		for i := 0; i < cfg.Weight; i++ {
			out = append(out, s)
		}
	}
	return out, nil
}

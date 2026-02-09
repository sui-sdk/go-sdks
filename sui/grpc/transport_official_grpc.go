//go:build official_grpc

package grpc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/structpb"
)

const defaultCallMethodPath = "/sui.rpc.v2.Service/Call"

type OfficialGRPCTransportOptions struct {
	Target     string
	MethodPath string
	Timeout    time.Duration
	DialOpts   []grpc.DialOption
}

type OfficialGRPCTransport struct {
	conn       *grpc.ClientConn
	methodPath string
	timeout    time.Duration
}

func NewOfficialGRPCTransport(opts OfficialGRPCTransportOptions) (Transport, error) {
	methodPath := opts.MethodPath
	if methodPath == "" {
		methodPath = defaultCallMethodPath
	}
	timeout := opts.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	target, useInsecure, err := normalizeGRPCTarget(opts.Target)
	if err != nil {
		return nil, err
	}

	dialOpts := make([]grpc.DialOption, 0, len(opts.DialOpts)+1)
	dialOpts = append(dialOpts, opts.DialOpts...)
	if useInsecure {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})))
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, target, dialOpts...)
	if err != nil {
		return nil, fmt.Errorf("grpc dial %q failed: %w", target, err)
	}

	return &OfficialGRPCTransport{conn: conn, methodPath: methodPath, timeout: timeout}, nil
}

func (t *OfficialGRPCTransport) Call(ctx context.Context, method string, params []any, out any) error {
	payload := map[string]any{"method": method, "params": params}
	req, err := structpb.NewStruct(payload)
	if err != nil {
		return fmt.Errorf("build grpc request failed: %w", err)
	}

	callCtx := ctx
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		callCtx, cancel = context.WithTimeout(ctx, t.timeout)
		defer cancel()
	}

	var resp structpb.Struct
	if err := t.conn.Invoke(callCtx, t.methodPath, req, &resp); err != nil {
		return fmt.Errorf("grpc invoke failed: %w", err)
	}

	respMap := resp.AsMap()
	if rpcErr, hasErr := respMap["error"]; hasErr && rpcErr != nil {
		return fmt.Errorf("grpc response error: %v", rpcErr)
	}

	result, ok := respMap["result"]
	if !ok {
		result = respMap
	}

	if out == nil {
		return nil
	}
	b, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshal grpc result failed: %w", err)
	}
	if err := json.Unmarshal(b, out); err != nil {
		return fmt.Errorf("decode grpc result failed: %w", err)
	}
	return nil
}

func (t *OfficialGRPCTransport) Close() error {
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

func normalizeGRPCTarget(raw string) (target string, useInsecure bool, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", false, fmt.Errorf("grpc target is required")
	}

	if strings.Contains(raw, "://") {
		u, err := url.Parse(raw)
		if err != nil {
			return "", false, fmt.Errorf("invalid grpc target %q: %w", raw, err)
		}
		switch u.Scheme {
		case "http":
			return u.Host, true, nil
		case "https":
			return u.Host, false, nil
		default:
			if u.Host != "" {
				return u.Host, false, nil
			}
			return raw, false, nil
		}
	}

	if strings.HasPrefix(raw, "127.") || strings.HasPrefix(raw, "localhost") {
		return raw, true, nil
	}
	return raw, false, nil
}

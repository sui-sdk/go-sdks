//go:build !official_grpc

package grpc

import (
	"fmt"
	"time"
)

type OfficialGRPCTransportOptions struct {
	Target     string
	MethodPath string
	Timeout    time.Duration
	DialOpts   []any
}

func NewOfficialGRPCTransport(opts OfficialGRPCTransportOptions) (Transport, error) {
	_ = opts
	return nil, fmt.Errorf("official grpc transport requires build tag 'official_grpc' and google.golang.org/grpc dependencies")
}

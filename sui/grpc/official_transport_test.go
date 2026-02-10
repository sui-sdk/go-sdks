package grpc

import (
	"context"
	"net"
	"testing"

	gogrpc "google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestOfficialGRPCTransportAndClient(t *testing.T) {
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer lis.Close()

	srv := gogrpc.NewServer(gogrpc.UnknownServiceHandler(func(_ any, stream gogrpc.ServerStream) error {
		req := &structpb.Struct{}
		if err := stream.RecvMsg(req); err != nil {
			return err
		}
		in := req.AsMap()
		method, _ := in["method"].(string)

		var result any
		switch method {
		case "suix_getReferenceGasPrice":
			result = "1000"
		case "sui_getObject":
			result = map[string]any{"objectId": "0x1"}
		case "sui_multiGetObjects":
			result = []any{map[string]any{"objectId": "0x1"}}
		default:
			result = map[string]any{}
		}

		out, err := structpb.NewStruct(map[string]any{"result": result})
		if err != nil {
			return err
		}
		return stream.SendMsg(out)
	}))
	defer srv.Stop()

	go func() { _ = srv.Serve(lis) }()

	client, err := NewClient(ClientOptions{
		Network: "localnet",
		BaseURL: "http://" + lis.Addr().String(),
	})
	if err != nil {
		t.Fatalf("new client failed: %v", err)
	}
	defer func() { _ = client.Close() }()

	if gas, err := client.GetReferenceGasPrice(context.Background()); err != nil || gas != "1000" {
		t.Fatalf("get gas price failed, gas=%q err=%v", gas, err)
	}

	objectID := "0x0000000000000000000000000000000000000000000000000000000000000001"
	if got, err := client.GetObject(context.Background(), objectID, nil); err != nil {
		t.Fatalf("get object failed: %v", err)
	} else if got["object"] == nil {
		t.Fatalf("expected wrapped object in response")
	}

	if got, err := client.GetObjects(context.Background(), []string{objectID}, nil); err != nil {
		t.Fatalf("get objects failed: %v", err)
	} else if got["objects"] == nil {
		t.Fatalf("expected wrapped objects in response")
	}
}

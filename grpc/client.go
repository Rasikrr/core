package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(ctx context.Context, addr string, opts ...grpc.DialOption) (*Client, error) {
	return NewClientWithOptions(ctx, addr, opts...)
}

func NewClientWithOptions(_ context.Context, addr string, opts ...grpc.DialOption) (*Client, error) {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, fmt.Errorf("create gRPC client dial connection: %w", err)
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Invoke(ctx context.Context, method string, args, reply any, opts ...grpc.CallOption) error {
	return c.conn.Invoke(ctx, method, args, reply, opts...)
}

func (c *Client) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return c.conn.NewStream(ctx, desc, method, opts...)
}

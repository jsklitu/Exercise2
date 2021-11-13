package main

import (
	proto "Exercise2/grpc"
	"context"
)

type Client struct {
	proto.UnimplementedCriticalSectionServiceServer
	id        int
	ipAndPort string
	timeStamp int
	queue     []string
}

func main() {

}

func (c *Client) Request(ctx context.Context, in *proto.Message) (*proto.Close, error) {
	return &proto.Close{}, nil
}

func (c *Client) Reply(ctx context.Context, in *proto.Message) (*proto.Close, error) {
	return &proto.Close{}, nil
}

func (c *Client) Release(ctx context.Context, in *proto.Message) (*proto.Close, error) {
	return &proto.Close{}, nil
}

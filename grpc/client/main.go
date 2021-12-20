package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	proto "tutorial/pkg/addservice"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := proto.NewAddServiceClient(conn)

	ctx := context.Background()
	req := &proto.Request{A: int64(10), B: int64(25)}
	if response, err := client.Add(ctx, req); err == nil {
		fmt.Println(response.Result)
	}

	if response, err := client.Multiply(ctx, req); err == nil {
		fmt.Println(response.Result)
	}

}

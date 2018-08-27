package main

import (
	"context"
	"fmt"
	"log"

	"github.com/athenabjorg/grpc/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer conn.Close()

	c := calculatorpb.NewCalculateServiceClient(conn)

	doUnary(c)
}

func doUnary(c calculatorpb.CalculateServiceClient) {
	fmt.Println("Running unary RPC...")
	req := &calculatorpb.CalculateRequest{
		Calculation: &calculatorpb.Calculation{
			A: 5,
			B: 7,
		},
	}

	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator RPC: %v", err)
	}

	log.Printf("Response from Calculator: %v", res.Result)
}

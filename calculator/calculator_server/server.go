package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/athenabjorg/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculateRequest) (*calculatorpb.CalculateResponse, error) {
	fmt.Printf("Calculatefunciton was invoked with %v", req)
	result := req.GetCalculation().GetA() + req.GetCalculation().GetB()
	res := &calculatorpb.CalculateResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	fmt.Println("Hello I'm a server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to lestan: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculateServiceServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

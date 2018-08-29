package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/athenabjorg/grpc/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum funciton was invoked with %v", req)
	result := req.GetValues().GetA() + req.GetValues().GetB()
	res := &calculatorpb.SumResponse{
		Result: result,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition funciton was invoked with %v", req)
	k := 2
	n := int(req.GetValue())
	for n > 1 {
		if n%k == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: strconv.Itoa(k),
			}
			stream.Send(res)
			n = n / k
		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage funciton was invoked with a streaming request\n")
	result := 0.0
	count := 0.0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: fmt.Sprintf("%.2f", result/count),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		result += float64(req.GetValue())
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Printf("FindMaximum funciton was invoked with a streaming request\n")

	maximum := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		number := int(req.GetValue())
		if number > maximum {
			maximum = number
			err = stream.Send(&calculatorpb.FindMaximumResponse{
				Result: strconv.Itoa(maximum),
			})
			if err != nil {
				log.Fatalf("Error while sending data to client: %v", err)
				return err
			}
		}
	}
}

func main() {
	fmt.Println("Hello I'm a server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	err = s.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
}

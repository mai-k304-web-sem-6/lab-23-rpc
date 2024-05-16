package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dariubs/percent"
	pb "github.com/mai-k304-web-sem-6/lab-23-rpc.git/pkg"
	"google.golang.org/grpc"
	"log"
	"math"
	"net"
	"strings"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement proto.GreeterServer.
type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) Sum(ctx context.Context, in *pb.TwoRequest) (*pb.Response, error) {
	return &pb.Response{Result: in.GetA() + in.GetB()}, nil
}

func (s *server) Subtract(ctx context.Context, in *pb.TwoRequest) (*pb.Response, error) {
	return &pb.Response{Result: in.GetA() - in.GetB()}, nil
}

func (s *server) Multiply(ctx context.Context, in *pb.TwoRequest) (*pb.Response, error) {
	return &pb.Response{Result: in.GetA() * in.GetB()}, nil
}

func (s *server) Share(ctx context.Context, in *pb.TwoRequest) (*pb.Response, error) {
	return &pb.Response{Result: in.GetA() / in.GetB()}, nil
}

func (s *server) Sqrt(ctx context.Context, in *pb.OneRequest) (*pb.Response, error) {
	return &pb.Response{Result: float32(math.Sqrt(float64(in.GetA())))}, nil
}

func (s *server) Percent(ctx context.Context, in *pb.TwoRequest) (*pb.Response, error) {
	return &pb.Response{Result: float32(percent.Percent(int(in.GetA()), int(in.GetB())))}, nil
}

func (s *server) Round(ctx context.Context, in *pb.TwoRequest) (*pb.Response, error) {
	return &pb.Response{Result: float32(math.Floor(float64(in.GetA())*math.Pow(10, float64(in.GetB()))+0.5) / math.Pow(10, float64(in.GetB())))}, nil
}

func (s *server) Exponentiation(ctx context.Context, in *pb.TwoRequest) (*pb.Response, error) {
	return &pb.Response{Result: float32(math.Pow(float64(in.GetA()), float64(in.GetB())))}, nil
}

func (s *server) Calculate(ctx context.Context, in *pb.CalculateRequest) (*pb.Response, error) {
	numbers := make([]float64, 0, len(in.GetNumbers()))
	for i := 0; i < len(in.GetNumbers()); i++ {
		numbers = append(numbers, float64(in.GetNumbers()[i]))
	}
	operations := strings.Split(in.GetOperations(), "")
	for i := 0; len(operations) > i; {
		switch operations[i] {
		case "2":
			res, _ := s.Sqrt(ctx, &pb.OneRequest{A: float32(numbers[i])})
			numbers[i] = float64(res.Result)
			operations = append(operations[:i], operations[i+1:]...)
		default:
			i++
		}
	}
	for i := 0; len(operations) > i; {
		switch operations[i] {
		case "%":
			res, _ := s.Percent(ctx, &pb.TwoRequest{A: float32(numbers[i]), B: float32(numbers[i+1])})
			numbers[i+1] = float64(res.Result)
			numbers = append(numbers[:i], numbers[i+1:]...)
			operations = append(operations[:i], operations[i+1:]...)
		case ":":
			res, _ := s.Round(ctx, &pb.TwoRequest{A: float32(numbers[i]), B: float32(numbers[i+1])})
			numbers[i+1] = float64(res.Result)
			numbers = append(numbers[:i], numbers[i+1:]...)
			operations = append(operations[:i], operations[i+1:]...)
		case "^":
			res, _ := s.Exponentiation(ctx, &pb.TwoRequest{A: float32(numbers[i]), B: float32(numbers[i+1])})
			numbers[i+1] = float64(res.Result)
			numbers = append(numbers[:i], numbers[i+1:]...)
			operations = append(operations[:i], operations[i+1:]...)
		default:
			i++
		}
	}
	for i := 0; len(operations) > i; {
		switch operations[i] {
		case "*":
			res, _ := s.Multiply(ctx, &pb.TwoRequest{A: float32(numbers[i]), B: float32(numbers[i+1])})
			numbers[i+1] = float64(res.Result)
			numbers = append(numbers[:i], numbers[i+1:]...)
			operations = append(operations[:i], operations[i+1:]...)
		case "/":
			res, _ := s.Share(ctx, &pb.TwoRequest{A: float32(numbers[i]), B: float32(numbers[i+1])})
			numbers[i+1] = float64(res.Result)
			numbers = append(numbers[:i], numbers[i+1:]...)
			operations = append(operations[:i], operations[i+1:]...)
		default:
			i++
		}
	}
	for i := 0; len(operations) > i; {
		switch operations[i] {
		case "+":
			res, _ := s.Sum(ctx, &pb.TwoRequest{A: float32(numbers[i]), B: float32(numbers[i+1])})
			numbers[i+1] = float64(res.Result)
			numbers = append(numbers[:i], numbers[i+1:]...)
			operations = append(operations[:i], operations[i+1:]...)
		case "-":
			res, _ := s.Subtract(ctx, &pb.TwoRequest{A: float32(numbers[i]), B: float32(numbers[i+1])})
			numbers[i+1] = float64(res.Result)
			numbers = append(numbers[:i], numbers[i+1:]...)
			operations = append(operations[:i], operations[i+1:]...)
		default:
			i++
		}
	}
	return &pb.Response{Result: float32(numbers[0])}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package server

import (
	"context"
	"log"

	pb "github.com/gwuah/postmates/api"
)

type GrpcServer struct {
	pb.UnimplementedPostmatesServer
}

func (s *GrpcServer) Signup(ctx context.Context, in *pb.SignupRequest) (*pb.SignupResponse, error) {
	log.Printf("Received: %v", in.GetPhone())
	return &pb.SignupResponse{Phone: in.GetPhone()}, nil
}

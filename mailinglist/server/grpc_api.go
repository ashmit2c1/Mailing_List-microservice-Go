package server

import (
	"context"
	"log"
	"mailinglist/db"
	"mailinglist/proto"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	DB *db.DB
	proto.UnimplementedMailingListServiceServer
}

func (s *GRPCServer) AddSubscriber(ctx context.Context, req *proto.Subscriber) (*proto.SubscriberResponse, error) {
	err := s.DB.AddSubscriber(req.Email)
	if err != nil {
		return &proto.SubscriberResponse{Success: false, Message: "Failed to add subscriber"}, nil
	}
	return &proto.SubscriberResponse{Success: true, Message: "Subscriber added"}, nil
}

func (s *GRPCServer) ListSubscribers(ctx context.Context, _ *proto.Empty) (*proto.SubscriberList, error) {
	emails, err := s.DB.ListSubscribers()
	if err != nil {
		return nil, err
	}

	var subscribers []*proto.Subscriber
	for _, email := range emails {
		subscribers = append(subscribers, &proto.Subscriber{Email: email})
	}
	return &proto.SubscriberList{Subscribers: subscribers}, nil
}

func StartGRPCServer(db *db.DB) {
	grpcServer := grpc.NewServer()
	proto.RegisterMailingListServiceServer(grpcServer, &GRPCServer{DB: db})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("Starting gRPC server on :50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

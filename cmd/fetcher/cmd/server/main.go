package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "content-parser/cmd/fetcher/proto"
	"content-parser/cmd/fetcher/service"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedFetcherServer
}

func (s *server) FetchAllRSS(ctx context.Context, in *pb.FetchAllRSSRequest) (*pb.FetchAllRSSReply, error) {
	c := service.NewClient()
	c.UpdateRSSFeeds()
	return &pb.FetchAllRSSReply{Message: "FetchAllRSS called successfully"}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterFetcherServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

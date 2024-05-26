package main

import (
	"log"
	"net"

	"github.com/Harikesh00/cloudbees/blog"
	"google.golang.org/grpc"
)

const portAddr = ":50051"

func main() {
	lis, err := net.Listen("tcp", portAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	blogService := blog.NewBlogServer()

	blog.RegisterBlogServiceServer(s, blogService)
	log.Printf("server listening on %s", portAddr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

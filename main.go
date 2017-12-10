package main

import (
	"flag"
	"log"

	"github.com/republicprotocol/republic/crypto"
	"github.com/republicprotocol/republic/dht"
	"github.com/republicprotocol/republic/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// Declare command line arguments.
var bind = flag.String("bind", "0.0.0.0", "Binding address")
var port = flag.String("port", "8080", "RPC listening port")

func main() {

	// Parse command line arguments.
	flag.Parse()

	// Generate identity for this node.
	secp, err := crypto.NewSECP256K1()
	if err != nil {
		log.Fatalf("failed to identify self: %v", err)
	}
	id := secp.PublicAddress()
	log.Println("Republic address:", id)

	// listen to the tcp port
	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	// Create gRPC services.
	node := NewNode(dht.ID(id))
	rpc.RegisterNodeServer(s, node)

	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Printf("Listening for connections on %s...\n", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
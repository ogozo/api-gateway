package client

import (
	"log"

	pb "github.com/ogozo/proto-definitions/gen/go/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitOrderServiceClient(orderServiceURL string) pb.OrderServiceClient {
	conn, err := grpc.NewClient(orderServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to order service: %v", err)
	}
	return pb.NewOrderServiceClient(conn)
}

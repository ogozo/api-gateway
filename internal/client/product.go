package client

import (
	"log"

	pb "github.com/ogozo/proto-definitions/gen/go/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitProductServiceClient(productServiceURL string) pb.ProductServiceClient {
	conn, err := grpc.NewClient(productServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to product service: %v", err)
	}
	return pb.NewProductServiceClient(conn)
}

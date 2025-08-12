package client

import (
	"log"

	pb "github.com/ogozo/proto-definitions/gen/go/cart"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitCartServiceClient(cartServiceURL string) pb.CartServiceClient {
	conn, err := grpc.NewClient(cartServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to cart service: %v", err)
	}
	return pb.NewCartServiceClient(conn)
}

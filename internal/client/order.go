package client

import (
	"log"

	pb "github.com/ogozo/proto-definitions/gen/go/order"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitOrderServiceClient(orderServiceURL string) pb.OrderServiceClient {
	conn, err := grpc.NewClient(orderServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}
	return pb.NewOrderServiceClient(conn)
}

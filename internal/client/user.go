package client

import (
	"log"

	pb "github.com/ogozo/proto-definitions/gen/go/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserServiceClient(userServiceURL string) pb.UserServiceClient {
	conn, err := grpc.NewClient(userServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}
	return pb.NewUserServiceClient(conn)
}

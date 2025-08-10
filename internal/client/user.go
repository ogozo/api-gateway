package client

import (
	"log"

	pb "github.com/ogozo/proto-definitions/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InitUserServiceClient, user servisine bir gRPC bağlantısı başlatır.
func InitUserServiceClient(userServiceURL string) pb.UserServiceClient {
	conn, err := grpc.Dial(userServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("could not connect to user service: %v", err)
	}

	return pb.NewUserServiceClient(conn)
}

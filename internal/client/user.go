package client

import (
	"context"

	"github.com/ogozo/api-gateway/internal/logging"
	pb "github.com/ogozo/proto-definitions/gen/go/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitUserServiceClient(userServiceURL string) pb.UserServiceClient {
	conn, err := grpc.NewClient(userServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logging.FromContext(context.Background()).Fatal("could not connect to user service",
			zap.Error(err),
			zap.String("url", userServiceURL),
		)
	}
	return pb.NewUserServiceClient(conn)
}

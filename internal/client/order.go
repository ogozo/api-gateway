package client

import (
	"context"

	"github.com/ogozo/api-gateway/internal/logging"
	pb "github.com/ogozo/proto-definitions/gen/go/order"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitOrderServiceClient(orderServiceURL string) pb.OrderServiceClient {
	conn, err := grpc.NewClient(orderServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logging.FromContext(context.Background()).Fatal("could not connect to order service",
			zap.Error(err),
			zap.String("url", orderServiceURL),
		)
	}
	return pb.NewOrderServiceClient(conn)
}

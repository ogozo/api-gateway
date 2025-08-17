package client

import (
	"context"

	"github.com/ogozo/api-gateway/internal/logging"
	pb "github.com/ogozo/proto-definitions/gen/go/cart"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitCartServiceClient(cartServiceURL string) pb.CartServiceClient {
	conn, err := grpc.NewClient(cartServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logging.FromContext(context.Background()).Fatal("could not connect to cart service",
			zap.Error(err),
			zap.String("url", cartServiceURL),
		)
	}
	return pb.NewCartServiceClient(conn)
}

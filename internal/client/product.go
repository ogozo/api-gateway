package client

import (
	"context"

	"github.com/ogozo/api-gateway/internal/logging"
	pb "github.com/ogozo/proto-definitions/gen/go/product"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitProductServiceClient(productServiceURL string) pb.ProductServiceClient {
	conn, err := grpc.NewClient(productServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		logging.FromContext(context.Background()).Fatal("could not connect to product service",
			zap.Error(err),
			zap.String("url", productServiceURL),
		)
	}
	return pb.NewProductServiceClient(conn)
}

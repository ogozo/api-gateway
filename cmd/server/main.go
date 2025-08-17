package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ogozo/api-gateway/internal/client"
	"github.com/ogozo/api-gateway/internal/config"
	"github.com/ogozo/api-gateway/internal/handler"
	"github.com/ogozo/api-gateway/internal/logging"
	"github.com/ogozo/api-gateway/internal/middleware"
	"github.com/ogozo/api-gateway/internal/observability"
	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var cfg config.GatewayConfig
	config.LoadConfig(&cfg)

	logging.Init(cfg.OtelServiceName)
	defer logging.Sync()

	api_logger := logging.FromContext(ctx)

	shutdown, err := observability.InitTracerProvider(ctx, cfg.OtelServiceName, cfg.OtelExporterEndpoint, api_logger)
	if err != nil {
		api_logger.Fatal("failed to initialize tracer provider", zap.Error(err))
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			api_logger.Fatal("failed to shutdown tracer provider", zap.Error(err))
		}
	}()

	userClient := client.InitUserServiceClient(cfg.UserServiceURL)
	productClient := client.InitProductServiceClient(cfg.ProductServiceURL)
	cartClient := client.InitCartServiceClient(cfg.CartServiceURL)
	orderClient := client.InitOrderServiceClient(cfg.OrderServiceURL)

	userHandler := handler.NewUserHandler(userClient)
	productHandler := handler.NewProductHandler(productClient)
	cartHandler := handler.NewCartHandler(cartClient)
	orderHandler := handler.NewOrderHandler(orderClient, cartClient, productClient)

	app := fiber.New()

	prometheus := fiberprometheus.New(cfg.OtelServiceName)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	app.Use(otelfiber.Middleware())
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	products := v1.Group("/products")
	products.Get("/:id", productHandler.GetProduct)
	products.Post("/", middleware.AuthRequired(cfg.JWTSecretKey), middleware.RoleRequired("ADMIN"), productHandler.CreateProduct)

	cart := v1.Group("/cart")
	cart.Use(middleware.AuthRequired(cfg.JWTSecretKey))
	cart.Get("/", cartHandler.GetCart)
	cart.Post("/items", cartHandler.AddItemToCart)

	v1.Post("/checkout", middleware.AuthRequired(cfg.JWTSecretKey), orderHandler.Checkout)

	protected := v1.Group("/me")
	protected.Use(middleware.AuthRequired(cfg.JWTSecretKey))
	protected.Get("/profile", userHandler.GetProfile)

	admin := v1.Group("/admin")
	admin.Use(middleware.AuthRequired(cfg.JWTSecretKey), middleware.RoleRequired("ADMIN"))
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the admin dashboard!"})
	})

	api_logger.Info("starting API Gateway", zap.String("port", cfg.HTTPPort))
	if err := app.Listen(cfg.HTTPPort); err != nil {
		api_logger.Fatal("failed to start server", zap.Error(err))
	}
}

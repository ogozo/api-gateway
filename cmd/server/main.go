package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ogozo/api-gateway/internal/client"
	"github.com/ogozo/api-gateway/internal/config"
	"github.com/ogozo/api-gateway/internal/handler"
	"github.com/ogozo/api-gateway/internal/middleware"
)

func main() {
	config.LoadConfig()
	cfg := config.AppConfig

	// gRPC istemcilerini başlat
	userClient := client.InitUserServiceClient(cfg.UserServiceURL)
	productClient := client.InitProductServiceClient(cfg.ProductServiceURL)
	cartClient := client.InitCartServiceClient(cfg.CartServiceURL)
	orderClient := client.InitOrderServiceClient(cfg.OrderServiceURL)

	// HTTP handler'larını başlat
	userHandler := handler.NewUserHandler(userClient)
	productHandler := handler.NewProductHandler(productClient)
	cartHandler := handler.NewCartHandler(cartClient)
	orderHandler := handler.NewOrderHandler(orderClient, cartClient, productClient)

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// --- Route Tanımlamaları ---
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	products := v1.Group("/products")
	products.Get("/:id", productHandler.GetProduct)
	products.Post("/", middleware.AuthRequired(cfg.JWTSecretKey), middleware.RoleRequired("ADMIN"), productHandler.CreateProduct)

	cart := v1.Group("/cart")
	cart.Use(middleware.AuthRequired(cfg.JWTSecretKey))
	cart.Get("/", cartHandler.GetCart)
	cart.Post("/items", cartHandler.AddItemToCart)

	// CHECKOUT ROUTE'U
	v1.Post("/checkout", middleware.AuthRequired(cfg.JWTSecretKey), orderHandler.Checkout)

	protected := v1.Group("/me")
	protected.Use(middleware.AuthRequired(cfg.JWTSecretKey))
	protected.Get("/profile", userHandler.GetProfile)

	admin := v1.Group("/admin")
	admin.Use(middleware.AuthRequired(cfg.JWTSecretKey), middleware.RoleRequired("ADMIN"))
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the admin dashboard!"})
	})

	log.Printf("API Gateway listening on port %s", cfg.HTTPPort)
	if err := app.Listen(cfg.HTTPPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

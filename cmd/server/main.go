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
	// 1. Yapılandırmayı yükle
	config.LoadConfig()
	cfg := config.AppConfig

	// 2. gRPC istemcilerini yapılandırmadan alarak başlat
	userClient := client.InitUserServiceClient(cfg.UserServiceURL)
	productClient := client.InitProductServiceClient(cfg.ProductServiceURL)
	cartClient := client.InitCartServiceClient(cfg.CartServiceURL)

	// 3. HTTP handler'larını başlat
	userHandler := handler.NewUserHandler(userClient)
	productHandler := handler.NewProductHandler(productClient)
	cartHandler := handler.NewCartHandler(cartClient)

	// 4. Fiber uygulamasını başlat
	app := fiber.New()
	app.Use(logger.New())

	// 5. Route'ları tanımla
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public Routes
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	// Product Routes
	products := v1.Group("/products")
	products.Get("/:id", productHandler.GetProduct)
	// JWT anahtarını config'den AuthRequired middleware'ine parametre olarak geçiyoruz
	products.Post("/", middleware.AuthRequired(cfg.JWTSecretKey), middleware.RoleRequired("ADMIN"), productHandler.CreateProduct)

	// Cart Routes
	cart := v1.Group("/cart")
	cart.Use(middleware.AuthRequired(cfg.JWTSecretKey))
	cart.Get("/", cartHandler.GetCart)
	cart.Post("/items", cartHandler.AddItemToCart)

	// Protected User Route
	protected := v1.Group("/me")
	protected.Use(middleware.AuthRequired(cfg.JWTSecretKey))
	protected.Get("/profile", userHandler.GetProfile)

	// Admin-only Route
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthRequired(cfg.JWTSecretKey), middleware.RoleRequired("ADMIN"))
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the admin dashboard!"})
	})

	// 6. Sunucuyu yapılandırmadan aldığı port ile başlat
	log.Printf("API Gateway listening on port %s", cfg.HTTPPort)
	err := app.Listen(cfg.HTTPPort)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

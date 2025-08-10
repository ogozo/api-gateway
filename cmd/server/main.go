package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/ogozo/api-gateway/internal/client"
	"github.com/ogozo/api-gateway/internal/handler"
	"github.com/ogozo/api-gateway/internal/middleware"
)

const (
	userServiceURL    = "localhost:50051"
	productServiceURL = "localhost:50052"
	cartServiceURL    = "localhost:50053"
	httpPort          = ":3000"
)

func main() {
	// gRPC istemcilerini başlat
	userClient := client.InitUserServiceClient(userServiceURL)
	productClient := client.InitProductServiceClient(productServiceURL)
	cartClient := client.InitCartServiceClient(cartServiceURL)

	// HTTP handler'larını başlat
	userHandler := handler.NewUserHandler(userClient)
	productHandler := handler.NewProductHandler(productClient)
	cartHandler := handler.NewCartHandler(cartClient)

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	// User Routes
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	// Product Routes
	products := v1.Group("/products")
	products.Get("/:id", productHandler.GetProduct)
	products.Post("/", middleware.AuthRequired(), middleware.RoleRequired("ADMIN"), productHandler.CreateProduct)

	// Cart Routes
	cart := v1.Group("/cart")
	cart.Use(middleware.AuthRequired())
	cart.Get("/", cartHandler.GetCart)
	cart.Post("/items", cartHandler.AddItemToCart)

	protected := v1.Group("/me")
	protected.Use(middleware.AuthRequired())
	protected.Get("/profile", userHandler.GetProfile)

	admin := v1.Group("/admin")
	admin.Use(middleware.AuthRequired(), middleware.RoleRequired("ADMIN"))
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the admin dashboard!"})
	})

	log.Printf("API Gateway listening on port %s", httpPort)
	if err := app.Listen(httpPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

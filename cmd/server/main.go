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
	// Bu adresleri daha sonra config'den alacağız
	userServiceURL = "localhost:50051"
	httpPort       = ":3000"
)

func main() {
	// gRPC istemcisini başlat
	userClient := client.InitUserServiceClient(userServiceURL)

	// HTTP handler'larını başlat
	userHandler := handler.NewUserHandler(userClient)

	// Fiber (web server) uygulamasını başlat
	app := fiber.New()
	app.Use(logger.New()) // Gelen istekleri logla

	// Route'ları tanımla
	api := app.Group("/api")
	v1 := api.Group("/v1")

	// Public (herkese açık) route'lar
	v1.Post("/register", userHandler.Register)
	v1.Post("/login", userHandler.Login)

	// Protected (kimlik doğrulaması gerektiren) route'lar
	protected := v1.Group("/me")
	protected.Use(middleware.AuthRequired())
	protected.Get("/profile", userHandler.GetProfile)

	// Admin (sadece admin rolüyle erişilebilen) route örneği
	admin := v1.Group("/admin")
	admin.Use(middleware.AuthRequired(), middleware.RoleRequired("ADMIN"))
	admin.Get("/dashboard", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the admin dashboard!"})
	})

	log.Printf("API Gateway listening on port %s", httpPort)
	err := app.Listen(httpPort)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

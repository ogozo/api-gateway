package handler

import (
	"github.com/gofiber/fiber/v2"
	pb "github.com/ogozo/proto-definitions/gen/go/user"
)

type UserHandler struct {
	userClient pb.UserServiceClient
}

func NewUserHandler(userClient pb.UserServiceClient) *UserHandler {
	return &UserHandler{userClient: userClient}
}

// Register, HTTP register isteğini alır ve gRPC'ye yönlendirir.
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req pb.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	res, err := h.userClient.Register(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

// Login, HTTP login isteğini alır ve gRPC'ye yönlendirir.
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req pb.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	res, err := h.userClient.Login(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res)
}

// GetProfile, korumalı bir endpoint örneğidir.
func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	// Middleware tarafından context'e eklenen kullanıcı bilgilerini alıyoruz.
	userID := c.Locals("user_id")
	role := c.Locals("role")

	return c.JSON(fiber.Map{
		"message": "This is a protected route",
		"user_id": userID,
		"role":    role,
	})
}

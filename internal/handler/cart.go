package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ogozo/api-gateway/internal/logging"
	pb "github.com/ogozo/proto-definitions/gen/go/cart"
)

type CartHandler struct {
	cartClient pb.CartServiceClient
}

func NewCartHandler(cartClient pb.CartServiceClient) *CartHandler {
	return &CartHandler{cartClient: cartClient}
}

func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	req := &pb.GetCartRequest{UserId: userID}

	res, err := h.cartClient.GetCart(c.UserContext(), req)
	if err != nil {
		logging.Error(c.UserContext(), "grpc call to cartClient.GetCart failed", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

func (h *CartHandler) AddItemToCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	type AddItemRequest struct {
		ProductID string `json:"product_id"`
		Quantity  int32  `json:"quantity"`
	}

	var reqBody AddItemRequest
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	req := &pb.AddItemToCartRequest{
		UserId: userID,
		Item: &pb.CartItem{
			ProductId: reqBody.ProductID,
			Quantity:  reqBody.Quantity,
		},
	}

	res, err := h.cartClient.AddItemToCart(c.UserContext(), req)
	if err != nil {
		logging.Error(c.UserContext(), "grpc call to cartClient.AddItemToCart failed", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

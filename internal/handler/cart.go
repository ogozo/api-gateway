package handler

import (
	"github.com/gofiber/fiber/v2"
	pb "github.com/ogozo/proto-definitions/gen/go/cart"
)

type CartHandler struct {
	cartClient pb.CartServiceClient
}

func NewCartHandler(cartClient pb.CartServiceClient) *CartHandler {
	return &CartHandler{cartClient: cartClient}
}

// GetCart, kullanıcının sepetini getirir.
func (h *CartHandler) GetCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string) // Auth middleware'inden gelir
	req := &pb.GetCartRequest{UserId: userID}

	res, err := h.cartClient.GetCart(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

// AddItemToCart, kullanıcının sepetine ürün ekler.
func (h *CartHandler) AddItemToCart(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string) // Auth middleware'inden gelir

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

	res, err := h.cartClient.AddItemToCart(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(res)
}

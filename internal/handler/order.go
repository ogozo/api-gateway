package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ogozo/api-gateway/internal/logging"
	pb_cart "github.com/ogozo/proto-definitions/gen/go/cart"
	pb_order "github.com/ogozo/proto-definitions/gen/go/order"
	pb_prod "github.com/ogozo/proto-definitions/gen/go/product"
	"go.uber.org/zap"
)

type OrderHandler struct {
	orderClient   pb_order.OrderServiceClient
	cartClient    pb_cart.CartServiceClient
	productClient pb_prod.ProductServiceClient
}

func NewOrderHandler(orderClient pb_order.OrderServiceClient, cartClient pb_cart.CartServiceClient, productClient pb_prod.ProductServiceClient) *OrderHandler {
	return &OrderHandler{
		orderClient:   orderClient,
		cartClient:    cartClient,
		productClient: productClient,
	}
}

func (h *OrderHandler) Checkout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	cartRes, err := h.cartClient.GetCart(c.UserContext(), &pb_cart.GetCartRequest{UserId: userID})
	if err != nil {
		logging.Error(c.UserContext(), "failed to get cart for checkout", err, zap.String("user_id", userID))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve cart"})
	}
	if len(cartRes.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cart is empty"})
	}

	var orderItems []*pb_order.OrderItem
	for _, cartItem := range cartRes.Items {
		prodReq := &pb_prod.GetProductRequest{ProductId: cartItem.ProductId}
		prodRes, err := h.productClient.GetProduct(c.UserContext(), prodReq)
		if err != nil {
			logging.Error(c.UserContext(), "failed to validate product price for checkout", err, zap.String("product_id", cartItem.ProductId))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not validate product price"})
		}

		orderItems = append(orderItems, &pb_order.OrderItem{
			ProductId: cartItem.ProductId,
			Quantity:  cartItem.Quantity,
			Price:     prodRes.Product.Price,
		})
	}

	orderReq := &pb_order.CreateOrderRequest{
		UserId: userID,
		Items:  orderItems,
	}
	orderRes, err := h.orderClient.CreateOrder(c.UserContext(), orderReq)
	if err != nil {
		logging.Error(c.UserContext(), "failed to create order", err, zap.String("user_id", userID))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create order"})
	}

	return c.Status(fiber.StatusAccepted).JSON(orderRes.Order)
}

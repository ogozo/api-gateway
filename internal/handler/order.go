package handler

import (
	"log"

	pb_cart "github.com/ogozo/proto-definitions/gen/go/cart"
	pb_order "github.com/ogozo/proto-definitions/gen/go/order"
	pb_prod "github.com/ogozo/proto-definitions/gen/go/product"
	"github.com/gofiber/fiber/v2"
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

// Checkout, Saga akışını başlatır.
func (h *OrderHandler) Checkout(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	// Adım 1: Kullanıcının sepetini al.
	cartRes, err := h.cartClient.GetCart(c.Context(), &pb_cart.GetCartRequest{UserId: userID})
	if err != nil {
		log.Printf("Error getting cart for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve cart"})
	}
	if len(cartRes.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cart is empty"})
	}

	// Adım 2: Güvenlik ve Tutarlılık için Ürün Fiyatlarını Doğrula
	// Sepetteki ürünlerin güncel fiyatlarını service-product'tan alıp doğrulamak,
	// kullanıcının sepetine ürün ekledikten sonra fiyat değişirse yanlış bir fiyattan
	// sipariş vermesini engeller.
	var orderItems []*pb_order.OrderItem
	for _, cartItem := range cartRes.Items {
		prodReq := &pb_prod.GetProductRequest{ProductId: cartItem.ProductId}
		prodRes, err := h.productClient.GetProduct(c.Context(), prodReq)
		if err != nil {
			log.Printf("Error validating product price for product %s: %v", cartItem.ProductId, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not validate product price"})
		}
		
		orderItems = append(orderItems, &pb_order.OrderItem{
			ProductId: cartItem.ProductId,
			Quantity:  cartItem.Quantity,
			Price:     prodRes.Product.Price, // Güvenli: Her zaman en güncel fiyatı alıyoruz.
		})
	}

	// 3. Adım: Sipariş oluşturma isteğini service-order'a gönder.
	orderReq := &pb_order.CreateOrderRequest{
		UserId: userID,
		Items:  orderItems,
	}
	orderRes, err := h.orderClient.CreateOrder(c.Context(), orderReq)
	if err != nil {
		log.Printf("Error creating order for user %s: %v", userID, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not create order"})
	}

	// Kullanıcıya siparişin işleme alındığını ve PENDING durumdaki halini döndür.
	// Kullanıcı, siparişin nihai durumunu (CONFIRMED/CANCELLED) daha sonra
	// GetOrder endpoint'inden sorgulayabilir.
	return c.Status(fiber.StatusAccepted).JSON(orderRes.Order)
}
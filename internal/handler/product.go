package handler

import (
	"github.com/gofiber/fiber/v2"
	pb "github.com/ogozo/proto-definitions/gen/go/product"
)

type ProductHandler struct {
	productClient pb.ProductServiceClient
}

func NewProductHandler(productClient pb.ProductServiceClient) *ProductHandler {
	return &ProductHandler{productClient: productClient}
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	req := &pb.GetProductRequest{ProductId: id}

	res, err := h.productClient.GetProduct(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res.Product)
}

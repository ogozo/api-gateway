package handler

import (
	pb "github.com/ogozo/proto-definitions/gen/go/product"
	"github.com/gofiber/fiber/v2"
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

	res, err := h.productClient.GetProduct(c.UserContext(), req)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(res.Product)
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req pb.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	res, err := h.productClient.CreateProduct(c.UserContext(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(res.Product)
}
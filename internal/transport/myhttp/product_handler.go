package myhttp

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/transport/myhttp/dto"
	"github.com/Vladislav5108/inventorytracker/internal/usecase"
)

type ProductHandler struct {
	catalog usecase.ProductCatalog
}

func NewProductHandler(catalog usecase.ProductCatalog) *ProductHandler {
	return &ProductHandler{catalog: catalog}
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil || id <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid product ID",
			"details": "ID must be a positive integer"})
		return
	}
	product, err := h.catalog.GetByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "product not found",
				"product": id,
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ProductResponse{
		ID:         product.ID,
		Name:       product.Name,
		Price:      product.Price,
		Quantity:   product.Quantity,
		CreatedAt:  product.CreatedAt,
		IsArchived: product.IsArchived,
	})
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	products, err := h.catalog.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			Price:      product.Price,
			Quantity:   product.Quantity,
			CreatedAt:  product.CreatedAt,
			IsArchived: product.IsArchived,
		})

	}
	c.JSON(http.StatusOK, response)
}
func (h *ProductHandler) GetByCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID must be an integer"})
		return
	}
	products, err := h.catalog.GetByCategory(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	response := make([]dto.ProductResponse, 0, len(products))
	for _, product := range products {
		response = append(response, dto.ProductResponse{
			ID:         product.ID,
			Name:       product.Name,
			Price:      product.Price,
			Quantity:   product.Quantity,
			CreatedAt:  product.CreatedAt,
			IsArchived: product.IsArchived,
		})
	}
	c.JSON(http.StatusOK, response)
}

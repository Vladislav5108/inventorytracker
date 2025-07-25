package myhttp

import (
	"errors"
	"net/http"
	"strconv"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/transport/myhttp/dto"
	"github.com/Vladislav5108/inventorytracker/internal/usecase"
	"github.com/gin-gonic/gin"
)

type AdminProductHandler struct {
	admin usecase.ProductAdmin
}

func NewAdminProductHandler(admin usecase.ProductAdmin) *AdminProductHandler {
	return &AdminProductHandler{admin: admin}
}

func (h *AdminProductHandler) Add(c *gin.Context) {
	var request dto.CreateProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format error"})
		return
	}
	product := domain.Product{
		Name:       request.Name,
		Price:      request.Price,
		Quantity:   request.Quantity,
		IsArchived: false,
	}
	id, err := h.admin.Add(product)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "add product"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *AdminProductHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive number "})
		return
	}

	var request dto.CreateProductRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect data format"})
		return
	}
	product := domain.Product{
		ID:       id,
		Name:     request.Name,
		Price:    request.Price,
		Quantity: request.Quantity,
	}
	err = h.admin.Update(product)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if errors.Is(err, domain.ErrDuplicateName) {
			c.JSON(http.StatusConflict, gin.H{"error": "duplicate name"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "product updated"})
}
func (h *AdminProductHandler) Archiv(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive number "})
		return
	}
	err = h.admin.Archiv(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on archiving"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":      "product archived",
		"is_archived": true})
}
func (h *AdminProductHandler) Restore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be integer"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive"})
		return
	}
	err = h.admin.Restore(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error on restore"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":     "product restore",
		"is_arhived": false,
	})
}
func (h *AdminProductHandler) GetArchived(c *gin.Context) {
	products, err := h.admin.GetArchived()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	response := make([]dto.ProductResponse, 0, len(products))
	for _, p := range products {
		response = append(response, dto.ProductResponse{
			ID:         p.ID,
			Name:       p.Name,
			Price:      p.Price,
			Quantity:   p.Quantity,
			CategoryID: p.CategoryID,
			IsArchived: true,
		})
	}

	c.JSON(200, response)
}

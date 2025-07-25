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

type CategoryHandler struct {
	category usecase.CategoryCatalog
}

func NewCategoryHandler(category usecase.CategoryCatalog) *CategoryHandler {
	return &CategoryHandler{category: category}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var request dto.CreateCategoryRequest

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}
	category := domain.Category{
		Name: request.Name,
	}
	id, err := h.category.CreateCategory(category)
	if err != nil {
		if errors.Is(err, domain.ErrAlreadyCategory) {
			c.JSON(http.StatusConflict, gin.H{"error": "category already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}
func (h *CategoryHandler) GetByIDCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}
	category, err := h.category.GetByIDCategory(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFoundCategory) {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, category)
}
func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	categories, err := h.category.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get categories"})
		return
	}
	if len(categories) == 0 {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}
	c.JSON(http.StatusOK, categories)
}
func (h *CategoryHandler) UpDateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	var request dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect data format"})
		return
	}
	category := domain.Category{
		ID:   id,
		Name: request.Name,
	}

	err = h.category.UpDateCategory(category)
	if err != nil {
		if errors.Is(err, domain.ErrNameCategory) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, category)
}
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}
	if id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}
	err = h.category.DeleteCategory(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFoundCategory) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found category"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, nil)
}

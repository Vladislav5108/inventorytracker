package myhttp_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/transport/myhttp"
	"github.com/Vladislav5108/inventorytracker/internal/transport/myhttp/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockProductCatalog struct {
	mock.Mock
}

func (m *MockProductCatalog) GetByID(id int) (domain.Product, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Product), args.Error(1)
}
func (m *MockProductCatalog) GetAll() ([]domain.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}
func (m *MockProductCatalog) GetByCategory(categoryID int) ([]domain.Product, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

var (
	testTime = time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	testProduct = domain.Product{
		ID:         1,
		Name:       "Test Product",
		Price:      100050,
		Quantity:   10,
		CategoryID: 1,
		IsArchived: false,
		CreatedAt:  testTime,
	}

	testProducts = []domain.Product{
		testProduct,
		{
			ID:         2,
			Name:       "Second Product",
			Price:      200075,
			Quantity:   5,
			CategoryID: 2,
			IsArchived: true,
			CreatedAt:  testTime.Add(24 * time.Hour),
		},
	}
)

func TestGetByID_Success(t *testing.T) {

	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	mockCatalog.On("GetByID", 1).Return(testProduct, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}

	handler.GetByID(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response dto.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, testProduct.ID, response.ID)
	assert.Equal(t, testProduct.Name, response.Name)
	assert.Equal(t, testProduct.Price, response.Price)
	assert.Equal(t, testProduct.Quantity, response.Quantity)
	assert.Equal(t, testProduct.CategoryID, response.CategoryID)
	assert.Equal(t, testProduct.IsArchived, response.IsArchived)

	mockCatalog.AssertExpectations(t)
}

func TestGetByID_InvalidID(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	t.Run("Non-numeric ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "abc"}}

		handler.GetByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		assert.Equal(t, "Invalid product ID", response["error"])
		assert.Equal(t, "ID must be a positive integer", response["details"])
	})

	t.Run("Negative ID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "-1"}}

		handler.GetByID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	mockCatalog.AssertNotCalled(t, "GetByID")
}
func TestGetByID_NotFound(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	nonExistentID := 999
	mockCatalog.On("GetByID", nonExistentID).Return(domain.Product{}, domain.ErrProductNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(nonExistentID)}}

	handler.GetByID(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "product not found", response["error"])
	assert.Equal(t, float64(nonExistentID), response["product"])

	mockCatalog.AssertExpectations(t)
}
func TestGetByID_ServerError(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	testID := 1
	mockCatalog.On("GetByID", testID).Return(domain.Product{}, errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(testID)}}

	handler.GetByID(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "internal server error", response["error"])
	mockCatalog.AssertExpectations(t)
}
func TestGetAll_Success(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	mockCatalog.On("GetAll").Return(testProducts, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 2, len(response))
}

func TestGetAll_EmptyList(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	mockCatalog.On("GetAll").Return([]domain.Product{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler.GetAll(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Empty(t, response)
	mockCatalog.AssertExpectations(t)
}

func TestGetAll_ServerError(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	expectedError := errors.New("database connection failed")
	mockCatalog.On("GetAll").Return(nil, expectedError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler.GetAll(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedError.Error(), response["error"])
	mockCatalog.AssertExpectations(t)
}
func TestGetAll_ContentType(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	mockCatalog.On("GetAll").Return(testProducts, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	handler.GetAll(c)

	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	mockCatalog.AssertExpectations(t)
}
func TestGetByCategory_Success(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	categoryID := 1
	expectedProducts := []domain.Product{
		{
			ID:         1,
			Name:       "Test Product",
			Price:      100050,
			Quantity:   10,
			CategoryID: 1,
			IsArchived: false,
			CreatedAt:  testTime,
		},
	}

	mockCatalog.On("GetByCategory", categoryID).Return(expectedProducts, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(categoryID)}}

	handler.GetByCategory(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, len(expectedProducts), len(response))
	assert.Equal(t, expectedProducts[0].ID, response[0].ID)
	assert.Equal(t, expectedProducts[0].Name, response[0].Name)
	assert.Equal(t, expectedProducts[0].CategoryID, response[0].CategoryID)

	mockCatalog.AssertCalled(t, "GetByCategory", categoryID)
	mockCatalog.AssertExpectations(t)
}

func TestGetByCategory_InvalidID(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	testCases := []struct {
		name     string
		idParam  string
		expected string
	}{
		{"Non-numeric ID", "abc", "ID must be an integer"},
		{"Empty ID", "", "ID must be an integer"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tc.idParam}}

			handler.GetByCategory(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			assert.Equal(t, tc.expected, response["error"])
		})
	}
	mockCatalog.AssertNotCalled(t, "GetByCategory")
}

func TestGetByCategory_NotFound(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	nonExistentCategoryID := 999
	mockCatalog.On("GetByCategory", nonExistentCategoryID).Return(nil, domain.ErrNotFoundCategory)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(nonExistentCategoryID)}}

	handler.GetByCategory(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, domain.ErrNotFoundCategory.Error(), response["error"])
	mockCatalog.AssertExpectations(t)
}

func TestGetByCategory_EmptyList(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	categoryID := 3
	mockCatalog.On("GetByCategory", categoryID).Return([]domain.Product{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(categoryID)}}

	handler.GetByCategory(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Empty(t, response)
	mockCatalog.AssertExpectations(t)
}

func TestGetByCategory_ServerError(t *testing.T) {
	mockCatalog := new(MockProductCatalog)
	handler := myhttp.NewProductHandler(mockCatalog)

	categoryID := 1
	expectedError := errors.New("database error")

	mockCatalog.On("GetByCategory", categoryID).Return(nil, expectedError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: strconv.Itoa(categoryID)}}

	handler.GetByCategory(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedError.Error(), response["error"])
	mockCatalog.AssertExpectations(t)
}

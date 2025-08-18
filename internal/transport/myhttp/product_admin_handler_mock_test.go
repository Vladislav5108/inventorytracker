package myhttp_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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

type MockProductAdmin struct {
	mock.Mock
}

func (m *MockProductAdmin) Add(product domain.Product) (int, error) {
	args := m.Called(product)
	return args.Int(0), args.Error(1)
}

func (m *MockProductAdmin) Update(product domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductAdmin) Archiv(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductAdmin) Restore(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockProductAdmin) GetArchived() ([]domain.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

func TestAdminProductHandler_Add_Success(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	requestBody := `{
        "name": "New Product",
        "price": 100050,
        "quantity": 10,
        "category_id": 1
    }`

	expectedProduct := domain.Product{
		Name:       "New Product",
		Price:      100050,
		Quantity:   10,
		CategoryID: 1,
		IsArchived: false,
	}

	mockAdmin.On("Add", expectedProduct).Return(123, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Add(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 123, response["id"])
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Add_InvalidJSON(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	invalidJSON := `{
	"name": "New Product",
	"price": "invalid",
	"quantity": 10,
	"category_id": 1
	}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(invalidJSON))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Add(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockAdmin.AssertNotCalled(t, "Add")
}

func TestAdminProductHandler_Add_MissingCategoryID(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	requestBody := `{
	"name": "New Product",
	"price": 100050,
	"quantity": 10
	}`

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Add(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockAdmin.AssertNotCalled(t, "Add")
}

func TestAdminProductHandler_Add_ForeignKeyError(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	requestBody := `{
	"name": "New Product",
        "price": 100050,
        "quantity": 10,
        "category_id": 999
	}`

	expectedProduct := domain.Product{
		Name:       "New Product",
		Price:      100050,
		Quantity:   10,
		CategoryID: 999,
		IsArchived: false,
	}

	mockAdmin.On("Add", expectedProduct).Return(0, errors.New("oreign key violation"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Add(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockAdmin.AssertExpectations(t)
}
func TestAdminProductHandler_Add_ServerError(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	requestBody := `{
        "name": "New Product",
        "price": 100050,
        "quantity": 10,
        "category_id": 1
    }`

	expectedProduct := domain.Product{
		Name:       "New Product",
		Price:      100050,
		Quantity:   10,
		CategoryID: 1,
		IsArchived: false,
	}

	mockAdmin.On("Add", expectedProduct).Return(0, errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Add(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Update_Success(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	requestBody := `{
        "name": "Updated Product",
        "price": 200075,
        "quantity": 5,
        "category_id": 2
    }`

	expectedProduct := domain.Product{
		ID:         1,
		Name:       "Updated Product",
		Price:      200075,
		Quantity:   5,
		CategoryID: 2,
		CreatedAt:  time.Time{},
		IsArchived: false,
	}

	mockAdmin.On("Update", mock.MatchedBy(func(p domain.Product) bool {
		return p.ID == expectedProduct.ID &&
			p.Name == expectedProduct.Name &&
			p.Price == expectedProduct.Price &&
			p.Quantity == expectedProduct.Quantity &&
			p.CategoryID == expectedProduct.CategoryID &&
			p.IsArchived == expectedProduct.IsArchived
	})).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "product updated", response["status"])
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Update_InvalidID(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	testCases := []struct {
		name    string
		idParam string
	}{
		{"Non-numeric ID", "abc"},
		{"Negative ID", "-1"},
		{"Zero ID", "0"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "id", Value: tc.idParam}}
			c.Request = httptest.NewRequest("PUT", "/", strings.NewReader(`{}`))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.Update(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			mockAdmin.AssertNotCalled(t, "Update")
		})
	}
}

func TestAdminProductHandler_Update_NotFound(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	requestBody := `{
        "name": "Non-existent",
        "price": 100050,
        "quantity": 10,
        "category_id": 1
    }`

	expectedProduct := domain.Product{
		ID:         999,
		Name:       "Non-existent",
		Price:      100050,
		Quantity:   10,
		CategoryID: 1,
	}

	mockAdmin.On("Update", expectedProduct).Return(domain.ErrProductNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "999"}}
	c.Request = httptest.NewRequest("PUT", "/", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Update_DuplicateName(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	requestBody := `{
        "name": "Duplicate",
        "price": 100050,
        "quantity": 10,
        "category_id": 1
    }`

	expectedProduct := domain.Product{
		ID:         1,
		Name:       "Duplicate",
		Price:      100050,
		Quantity:   10,
		CategoryID: 1,
	}

	mockAdmin.On("Update", expectedProduct).Return(domain.ErrDuplicateName)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusConflict, w.Code)
	mockAdmin.AssertExpectations(t)
}
func TestAdminProductHandler_Update_ServerError(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	requestBody := `{
        "name": "Test",
        "price": 100050,
        "quantity": 10,
        "category_id": 1
    }`

	expectedProduct := domain.Product{
		ID:         1,
		Name:       "Test",
		Price:      100050,
		Quantity:   10,
		CategoryID: 1,
	}

	mockAdmin.On("Update", expectedProduct).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/", strings.NewReader(requestBody))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Update(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Archive_Success(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)
	productID := 1

	mockAdmin.On("Archiv", productID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{Key: "id", Value: strconv.Itoa(productID)},
	}

	handler.Archiv(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "product archived", response["status"])
	assert.True(t, response["is_archived"].(bool))

	mockAdmin.AssertExpectations(t)
	mockAdmin.AssertCalled(t, "Archiv", productID)
}

func TestAdminProductHandler_Archive_InvalidID(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	testCases := []struct {
		name    string
		idParam string
	}{
		{"String ID", "abc"},
		{"Negative ID", "-1"},
		{"Zero ID", "0"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tc.idParam}}

			handler.Archiv(c)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			mockAdmin.AssertNotCalled(t, "Archive")
		})
	}
}

func TestAdminProductHandler_Archive_NotFound(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	nonExistentID := 999

	mockAdmin.On("Archiv", nonExistentID).Return(domain.ErrProductNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(nonExistentID)}}
	c.Request = httptest.NewRequest("PUT", "/", nil)

	handler.Archiv(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Archive_AlreadyArchived(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)
	productID := 1

	mockAdmin.On("Archiv", productID).Return(domain.ErrProductArchived)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(productID)}}
	c.Request = httptest.NewRequest("PUT", "/products/archive", nil)

	handler.Archiv(c)

	assert.Equal(t, http.StatusConflict, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "product on archive", response["error"])
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Archive_ServerError(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)
	productID := 1

	mockAdmin.On("Archiv", productID).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(productID)}}

	handler.Archiv(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Restore_Success(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	productID := 1

	mockAdmin.On("Restore", productID).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(productID)}}
	c.Request = httptest.NewRequest("PUT", "/", nil)

	handler.Restore(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "product restore", response["status"])
	assert.False(t, response["is_archived"].(bool))
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Restore_InvalidIDFormat(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "abc"}}
	c.Request = httptest.NewRequest("PUT", "/", nil)

	handler.Restore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "id must be integer", response["error"])
	mockAdmin.AssertNotCalled(t, "Restore")
}

func TestAdminProductHandler_Restore_NegativeID(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "-1"}}
	c.Request = httptest.NewRequest("PUT", "/", nil)

	handler.Restore(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "id must be a positive", response["error"])
	mockAdmin.AssertNotCalled(t, "Restore")
}

func TestAdminProductHandler_Restore_NotFound(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)
	nonExistentID := 999

	mockAdmin.On("Restore", nonExistentID).Return(domain.ErrProductNotFound)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(nonExistentID)}}
	c.Request = httptest.NewRequest("PUT", "/", nil)

	handler.Restore(c)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "product not found", response["error"])
	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_Restore_ServerError(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	productID := 1

	mockAdmin.On("Restore", productID).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: strconv.Itoa(productID)}}
	c.Request = httptest.NewRequest("PUT", "/", nil)

	handler.Restore(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "error on restore", response["error"])
	mockAdmin.AssertExpectations(t)
}
func TestAdminProductHandler_GetArchived_Success(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	testTime := time.Now()
	expectedProducts := []domain.Product{
		{
			ID:         1,
			Name:       "Archived Product 1",
			Price:      1000,
			Quantity:   5,
			CategoryID: 1,
			IsArchived: true,
			CreatedAt:  testTime,
		},
	}

	mockAdmin.On("GetArchived").Return(expectedProducts, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/archived", nil)

	handler.GetArchived(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, len(expectedProducts), len(response))

	if len(response) > 0 {
		assert.Equal(t, expectedProducts[0].ID, response[0].ID)
		assert.Equal(t, expectedProducts[0].Name, response[0].Name)
		assert.Equal(t, expectedProducts[0].IsArchived, response[0].IsArchived)
		assert.True(t, expectedProducts[0].CreatedAt.Equal(response[0].CreatedAt))
	}

	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_GetArchived_ServerError(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	mockAdmin.On("GetArchived").Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/archived", nil)

	handler.GetArchived(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "server error", response["error"])

	mockAdmin.AssertExpectations(t)
}

func TestAdminProductHandler_GetArchived_EmptyList(t *testing.T) {
	mockAdmin := new(MockProductAdmin)
	handler := myhttp.NewAdminProductHandler(mockAdmin)

	mockAdmin.On("GetArchived").Return([]domain.Product{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/archived", nil)

	handler.GetArchived(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.ProductResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Empty(t, response)

	mockAdmin.AssertExpectations(t)
}

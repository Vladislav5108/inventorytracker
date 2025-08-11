package myhttp_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/transport/myhttp"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCategoryCatalog struct {
	mock.Mock
}

func (m *MockCategoryCatalog) CreateCategory(category domain.Category) (int, error) {
	args := m.Called(category)
	return args.Int(0), args.Error(1)
}
func (m *MockCategoryCatalog) GetByIDCategory(id int) (domain.Category, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Category), args.Error(1)
}
func (m *MockCategoryCatalog) GetAllCategories() ([]domain.Category, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Category), args.Error(1)
}
func (m *MockCategoryCatalog) UpdateCategory(category domain.Category) error {
	args := m.Called(category)
	return args.Error(0)
}
func (m *MockCategoryCatalog) DeleteCategory(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter() (*gin.Engine, *MockCategoryCatalog) {
	mockCategory := new(MockCategoryCatalog)

	handler := myhttp.NewCategoryHandler(mockCategory)

	router := gin.Default()

	router.POST("/categories", handler.CreateCategory)
	router.GET("/categories/:id", handler.GetByIDCategory)
	router.GET("/categories", handler.GetAllCategories)
	router.PUT("/categories/:id", handler.UpDateCategory)
	router.DELETE("/categories/:id", handler.DeleteCategory)

	return router, mockCategory
}

func TestCreateCategory_Success(t *testing.T) {
	router, mockCategory := setupTestRouter()

	expectedCategory := domain.Category{
		Name:        "Test Category",
		Description: "Test Description",
	}

	mockCategory.On("CreateCategory", expectedCategory).Return(1, nil)

	requestBody := `{
	"name": "Test Category",
	"description": "Test Description"
	}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, 1, response["id"])

	mockCategory.AssertExpectations(t)
}

func TestCreateCategory_EmptyName(t *testing.T) {
	router, mockCategory := setupTestRouter()

	mockCategory.On("CreateCategory", domain.Category{
		Name:        "",
		Description: "Описание",
	}).Return(0, errors.New("name cannot be empty"))

	requestBody := `{"description": "Описание"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "internal server error", response["error"])

	mockCategory.AssertExpectations(t)
}

func TestCreateCategory_InvalidJSON(t *testing.T) {
	router, mockCategory := setupTestRouter()

	invalidJSON := `{
        "name": "Тестовая категория",
        "description": 12345
    }`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotEmpty(t, response["error"])

	mockCategory.AssertNotCalled(t, "CreateCategory")
}

func TestCreateCategory_EmptyBody(t *testing.T) {
	router, mockCategory := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockCategory.AssertNotCalled(t, "CreateCategory")
}
func TestCreateCategory_InvalidContentType(t *testing.T) {
	router, mockCategory := setupTestRouter()

	mockCategory.On("CreateCategory", domain.Category{
		Name:        "",
		Description: "",
	}).Return(1, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "text/plain")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]int
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, 1, response["id"])

	mockCategory.AssertExpectations(t)
}
func TestCreateCategory_ServiceError(t *testing.T) {
	router, mockCategory := setupTestRouter()

	mockCategory.On("CreateCategory", mock.Anything).Return(0, errors.New("db error"))

	requestBody := `{"name": "Test", "description": "Test"}`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/categories", bytes.NewBufferString(requestBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockCategory.AssertExpectations(t)
}

func TestGetByIDCategory_Success(t *testing.T) {
	router, mockCategory := setupTestRouter()

	expectedCategory := domain.Category{
		ID:   1,
		Name: "Test Category",
	}

	mockCategory.On("GetByIDCategory", 1).Return(expectedCategory, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, expectedCategory.ID, response.ID)
	assert.Equal(t, expectedCategory.Name, response.Name)

	mockCategory.AssertExpectations(t)
}

func TestGetByIDCategory_NotFound(t *testing.T) {
	router, mockCategory := setupTestRouter()

	mockCategory.On("GetByIDCategory", 999).Return(domain.Category{}, domain.ErrNotFoundCategory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "category not found", response["error"])
	mockCategory.AssertExpectations(t)
}

func TestGetByIDCategory_InvalidID_NotNumber(t *testing.T) {
	router, mockCategory := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "id must be an integer", response["error"])

	mockCategory.AssertNotCalled(t, "GetByIDCategory")
}

func TestGetByIDCategory_ServerError(t *testing.T) {
	router, mockCategory := setupTestRouter()

	mockCategory.On("GetByIDCategory", 1).Return(domain.Category{}, errors.New("database error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "internal server error", response["error"])
	mockCategory.AssertExpectations(t)
}

func TestGetByIDCategory_NegativeID(t *testing.T) {
	router, mockCategory := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories/-1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "id must be a positive integer", response["error"])

	mockCategory.AssertNotCalled(t, "GetByIDCategory")
}

func TestGetAllCategories_Success(t *testing.T) {

	router, mockCategory := setupTestRouter()

	expectedCategories := []domain.Category{
		{ID: 1, Name: "Category 1"},
		{ID: 2, Name: "Category 2"},
	}

	mockCategory.On("GetAllCategories").Return(expectedCategories, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []domain.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 2, len(response))

	assert.Equal(t, 1, response[0].ID)
	assert.Equal(t, "Category 1", response[0].Name)

	mockCategory.AssertExpectations(t)
}

func TestGetAllCategories_EmptyList(t *testing.T) {
	router, mockCategory := setupTestRouter()

	mockCategory.On("GetAllCategories").Return([]domain.Category{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 0, len(response))
	mockCategory.AssertExpectations(t)
}
func TestGetAllCategories_ServerError(t *testing.T) {
	router, mockCategory := setupTestRouter()

	mockCategory.On("GetAllCategories").Return(nil, errors.New("database error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/categories", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "failed to get categories", response["error"])
	mockCategory.AssertExpectations(t)
}

func TestUpDateCategory_Success(t *testing.T) {

	router, mockCategory := setupTestRouter()

	testID := 1
	requestBody := `{
        "name": "Updated Category",
        "description": "Updated Description"
    }`

	expectedCategory := domain.Category{
		ID:          testID,
		Name:        "Updated Category",
		Description: "Updated Description",
	}

	mockCategory.On("UpdateCategory", expectedCategory).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/categories/1", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response domain.Category
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, testID, response.ID)
	assert.Equal(t, "Updated Category", response.Name)
	assert.Equal(t, "Updated Description", response.Description)

	mockCategory.AssertExpectations(t)
}

func TestUpdateCategory_InvalidID_NotNumber(t *testing.T) {
	router, mockCategory := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/categories/abc", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "id must be an integer", response["error"])
	mockCategory.AssertNotCalled(t, "UpddateCategory")
}

func TestUpDateCategory_InvalidID_Negative(t *testing.T) {
	router, mockCategory := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/categories/-1", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "id must be a positive integer", response["error"])
	mockCategory.AssertNotCalled(t, "UpdateCategory")
}
func TestUpDateCategory_InvalidJSON(t *testing.T) {
	router, mockCategory := setupTestRouter()

	invalidJSON := `{
        "name": "Test",
        "description": 123
    }`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/categories/1", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "incorrect data format", response["error"])
	mockCategory.AssertNotCalled(t, "UpdateCategory")
}

func TestUpDateCategory_NotFound(t *testing.T) {
	router, mockCategory := setupTestRouter()

	testID := 999
	requestBody := `{
        "name": "Non-existent",
        "description": "Category"
    }`

	expectedCategory := domain.Category{
		ID:          testID,
		Name:        "Non-existent",
		Description: "Category",
	}

	mockCategory.On("UpdateCategory", expectedCategory).Return(domain.ErrNotFoundCategory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/categories/999", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "category not found", response["error"])
	mockCategory.AssertExpectations(t)
}

func TestUpDateCategory_ServerError(t *testing.T) {
	router, mockCategory := setupTestRouter()

	testID := 1
	requestBody := `{
        "name": "Test",
        "description": "Error case"
    }`

	expectedCategory := domain.Category{
		ID:          testID,
		Name:        "Test",
		Description: "Error case",
	}

	mockCategory.On("UpdateCategory", expectedCategory).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/categories/1", bytes.NewBufferString(requestBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "internal server error", response["error"])
	mockCategory.AssertExpectations(t)
}
func TestDeleteCategory_Success(t *testing.T) {
	router, mockCategory := setupTestRouter()

	testID := 1

	mockCategory.On("DeleteCategory", testID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "null", w.Body.String())

	mockCategory.AssertExpectations(t)
}

func TestDeleteCategory_InvalidID_NotNumber(t *testing.T) {
	router, mockCategory := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/abc", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "id must be an integer", response["error"])

	mockCategory.AssertNotCalled(t, "DeleteCategory")
}
func TestDeleteCategory_InvalidID_Negative(t *testing.T) {

	router, mockCategory := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/-1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "id must be a positive integer", response["error"])

	mockCategory.AssertNotCalled(t, "DeleteCategory")
}
func TestDeleteCategory_NotFound(t *testing.T) {
	router, mockCategory := setupTestRouter()

	testID := 999

	mockCategory.On("DeleteCategory", testID).Return(domain.ErrNotFoundCategory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "not found category", response["error"])

	mockCategory.AssertExpectations(t)
}
func TestDeleteCategory_ServerError(t *testing.T) {
	router, mockCategory := setupTestRouter()

	testID := 1

	mockCategory.On("DeleteCategory", testID).Return(errors.New("database error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/categories/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "failed to delete category", response["error"])

	mockCategory.AssertExpectations(t)
}

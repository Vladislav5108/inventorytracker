package usecase

import (
	"errors"
	"testing"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) CreateCategory(category domain.Category) (int, error) {
	args := m.Called(category)
	return args.Int(0), args.Error(1)
}
func (m *MockCategoryRepository) GetByIDCategory(id int) (domain.Category, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Category), args.Error(1)
}
func (m *MockCategoryRepository) UpdateCategory(category domain.Category) error {
	args := m.Called(category) // Фиксируем вызов и аргумент (category)
	return args.Error(0)       // Возвращаем только ошибку (метод ничего не возвращает кроме error)
}
func (m *MockCategoryRepository) DeleteCategory(id int) error {
	args := m.Called(id) // Фиксируем вызов и аргумент (id)
	return args.Error(0) // Возвращаем только ошибку
}
func (m *MockCategoryRepository) GetAllCategories() ([]domain.Category, error) {
	args := m.Called() // Фиксируем вызов (без аргументов)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Category), args.Error(1)
}

func TestCreateCategory_Success(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	input := domain.Category{Name: "Books"}
	repo.On("CreateCategory", input).Return(1, nil)

	id, err := uc.CreateCategory(input)

	assert.NoError(t, err)
	assert.Equal(t, 1, id)
	repo.AssertExpectations(t)
}
func TestCreateCategory_EmptyName(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	_, err := uc.CreateCategory(domain.Category{Name: ""})
	assert.ErrorIs(t, err, domain.ErrNameCategory)
	repo.AssertNotCalled(t, "CreateCategory")
}
func TestCreateCategory_WithID(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	_, err := uc.CreateCategory(domain.Category{ID: 10, Name: "Books"})

	assert.ErrorContains(t, err, "the id must not be specified")
	repo.AssertNotCalled(t, "CreateCategory")
}
func TestCreateCatgory_RepoError(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	input := domain.Category{Name: "Books"}
	expectedErr := errors.New("database error")

	repo.On("CreateCategory", input).Return(0, expectedErr)

	_, err := uc.CreateCategory(input)

	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}

func TestGetByIDCategory_InvalidID(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	_, err := uc.GetByIDCategory(-1)

	assert.ErrorIs(t, err, domain.ErrNotFoundCategory)
	repo.AssertNotCalled(t, "GetByIDCategory")
}

func TestGetByIDCategory_Success(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	expectedCategory := domain.Category{ID: 1, Name: "Books"}
	repo.On("GetByIDCategory", 1).Return(expectedCategory, nil)

	result, err := uc.GetByIDCategory(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedCategory, result)
	repo.AssertExpectations(t)
}
func TestGetByIDCategory_NotFound(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	repo.On("GetByIDCategory", 2).Return(domain.Category{}, domain.ErrNotFoundCategory)

	_, err := uc.GetByIDCategory(2)

	assert.ErrorContains(t, err, "repository error")
	assert.True(t, errors.Is(err, domain.ErrNotFoundCategory))
	repo.AssertExpectations(t)
}
func TestGetByIDCategory_RepositoryError(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	expectedErr := errors.New("database connection failed")
	repo.On("GetByIDCategory", 3).Return(domain.Category{}, expectedErr)

	_, err := uc.GetByIDCategory(3)

	assert.ErrorContains(t, err, "repository error")
	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}
func TestGetAllCategories_Success(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	expected := []domain.Category{
		{ID: 1, Name: "Books"},
		{ID: 2, Name: "Electronics"},
	}

	repo.On("GetAllCategories").Return(expected, nil)

	result, err := uc.GetAllCategories()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	repo.AssertExpectations(t)
}
func TestGetAllCategory_Empty(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	repo.On("GetAllCategories").Return([]domain.Category{}, nil)

	result, err := uc.GetAllCategories()

	assert.NoError(t, err)
	assert.Empty(t, result)
	repo.AssertExpectations(t)
}
func TestGetAllCategories_RepoError(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	expectedErr := errors.New("database connection failed")

	repo.On("GetAllCategories").Return(nil, expectedErr)

	result, err := uc.GetAllCategories()

	assert.ErrorContains(t, err, "error getting all categories")
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	repo.AssertExpectations(t)
}
func TestUpdateCategory_InvalidID(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	err := uc.UpdateCategory(domain.Category{ID: 0, Name: "Valid Name"})

	assert.ErrorIs(t, err, domain.ErrNotFoundCategory)
	repo.AssertNotCalled(t, "UpdateCategory")
}
func TestUpdateCategory_EmptyName(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	err := uc.UpdateCategory(domain.Category{ID: 1, Name: ""})

	assert.ErrorIs(t, err, domain.ErrNameCategory)
	repo.AssertNotCalled(t, "UpdateCategory")
}
func TestUpdateCategory_Success(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	testCategory := domain.Category{ID: 1, Name: "Updated Name"}
	repo.On("UpdateCategory", testCategory).Return(nil)

	err := uc.UpdateCategory(testCategory)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
func TestUpdateCategory_NotFound(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	testCategory := domain.Category{ID: 999, Name: "Valid"}
	expectedErr := errors.New("category not found")
	repo.On("UpdateCategory", testCategory).Return(expectedErr)

	err := uc.UpdateCategory(testCategory)

	assert.ErrorContains(t, err, "category update error")
	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}
func TestUpdateCategory_DBError(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	testCategory := domain.Category{ID: 1, Name: "Valid"}
	expectedErr := errors.New("database connection failed")
	repo.On("UpdateCategory", testCategory).Return(expectedErr)

	err := uc.UpdateCategory(testCategory)

	assert.ErrorContains(t, err, "category update error")
	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}
func TestDeleteCategory_InvalidID(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	err := uc.DeleteCategory(0)

	assert.ErrorIs(t, err, domain.ErrNotFoundCategory)
	repo.AssertNotCalled(t, "DeleteCategory")
}
func TestDeleteCategory_Success(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	testID := 1
	repo.On("DeleteCategory", testID).Return(nil)

	err := uc.DeleteCategory(testID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}
func TestDeleteCategory_NotFound(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	testID := 999
	expectedErr := errors.New("category not found")
	repo.On("DeleteCategory", testID).Return(expectedErr)

	err := uc.DeleteCategory(testID)

	assert.ErrorContains(t, err, "error delete category")
	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}
func TestDeleteCategory_DBError(t *testing.T) {
	repo := new(MockCategoryRepository)
	uc := NewCategoryUseCase(repo)

	testID := 1
	expectedErr := errors.New("database connection failed")
	repo.On("DeleteCategory", testID).Return(expectedErr)

	err := uc.DeleteCategory(testID)

	assert.ErrorContains(t, err, "error delete category")
	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}

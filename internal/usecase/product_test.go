package usecase

import (
	"errors"
	"testing"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetByID(id int) (domain.Product, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Product), args.Error(1)
}
func (m *MockProductRepository) GetAll() ([]domain.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}
func (m *MockProductRepository) GetByCategory(categoryID int) ([]domain.Product, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}
func TestGetByID_InvalidID(t *testing.T) {
	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	result, err := uc.GetByID(0)

	assert.Equal(t, domain.Product{}, result)
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
	repo.AssertNotCalled(t, "GetByID")
}
func TestGetByID_Success(t *testing.T) {
	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	expectedProduct := domain.Product{
		ID:   1,
		Name: "Телефон",
	}
	repo.On("GetByID", 1).Return(expectedProduct, nil)

	result, err := uc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, result)
	repo.AssertExpectations(t)
}
func TestGetByID_RepositoryError(t *testing.T) {
	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	expectedErr := errors.New("database error")

	repo.On("GetByID", 1).Return(domain.Product{}, expectedErr)

	result, err := uc.GetByID(1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Equal(t, domain.Product{}, result)

	repo.AssertExpectations(t)
}

func TestGetAll_Success(t *testing.T) {

	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	expectedProducts := []domain.Product{
		{ID: 1, Name: "Laptop"},
		{ID: 2, Name: "Phone"},
	}

	repo.On("GetAll").Return(expectedProducts, nil)

	products, err := uc.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	assert.Len(t, products, 2)
	repo.AssertExpectations(t)
}

func TestGetAll_RepositoryError(t *testing.T) {
	t.Parallel()

	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	expectedErr := errors.New("database error")
	repo.On("GetAll").Return(nil, expectedErr)

	products, err := uc.GetAll()

	require.Error(t, err)
	assert.Nil(t, products)
	assert.ErrorContains(t, err, "error getting products")
	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}

func TestGetAll_EmptyList(t *testing.T) {

	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	repo.On("GetAll").Return([]domain.Product{}, nil)

	products, err := uc.GetAll()

	assert.Error(t, err)
	assert.Nil(t, products)
	assert.ErrorContains(t, err, "the list is empty")
	repo.AssertExpectations(t)
}

func TestGetByCategory_InvalidID(t *testing.T) {
	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	products, err := uc.GetByCategory(0)

	require.Error(t, err)
	assert.Nil(t, products)
	assert.ErrorIs(t, err, domain.ErrNotFoundCategory)
	repo.AssertNotCalled(t, "GetByCategory")
}
func TestGetByCategory_Success(t *testing.T) {
	t.Parallel()

	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	categoryID := 5
	expectedProducts := []domain.Product{
		{ID: 1, Name: "Laptop", CategoryID: categoryID},
		{ID: 2, Name: "Mouse", CategoryID: categoryID},
	}

	repo.On("GetByCategory", categoryID).Return(expectedProducts, nil)

	products, err := uc.GetByCategory(categoryID)

	require.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	repo.AssertExpectations(t)
}
func TestGetByCategory_RepoError(t *testing.T) {
	t.Parallel()

	repo := new(MockProductRepository)
	uc := NewProductUseCase(repo)

	expectedErr := errors.New("database error")

	repo.On("GetByCategory", 3).Return(nil, expectedErr)

	products, err := uc.GetByCategory(3)

	require.Error(t, err)
	assert.Nil(t, products)
	assert.ErrorContains(t, err, "error getting product on category")
	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}

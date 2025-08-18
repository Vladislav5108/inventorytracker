package usecase

import (
	"errors"
	"testing"

	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockProductAdminRepository struct {
	mock.Mock
}

func (m *MockProductAdminRepository) Add(product domain.Product) (int, error) {
	args := m.Called(product)
	return args.Int(0), args.Error(1)
}
func (m *MockProductAdminRepository) Update(product domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *MockProductAdminRepository) Archiv(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockProductAdminRepository) Restore(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *MockProductAdminRepository) GetArchived() ([]domain.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

type MockCatalogRepository struct {
	mock.Mock
}

func (m *MockCatalogRepository) GetByID(id int) (domain.Product, error) {
	args := m.Called(id)
	return args.Get(0).(domain.Product), args.Error(1)
}
func (m *MockCatalogRepository) GetAll() ([]domain.Product, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

func (m *MockCatalogRepository) GetByCategory(categoryID int) ([]domain.Product, error) {
	args := m.Called(categoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Product), args.Error(1)
}

func TestAdminProductUseCase_Add_Success(t *testing.T) {
	repo := new(MockProductAdminRepository)
	uc := AdminProductUseCase{productRepo: repo}

	expectedProduct := domain.Product{
		Name:       "Ноутбук",
		Price:      1000,
		IsArchived: false,
	}
	repo.On("Add", expectedProduct).Return(1, nil)

	id, err := uc.Add(domain.Product{
		Name:  "Ноутбук",
		Price: 1000,
	})

	require.NoError(t, err)
	assert.Equal(t, 1, id)
	repo.AssertExpectations(t)
}
func TestAdminProductUseCase_Add_EmptyName(t *testing.T) {
	repo := new(MockProductAdminRepository)
	uc := AdminProductUseCase{productRepo: repo}

	id, err := uc.Add(domain.Product{
		Name:  "",
		Price: 1000,
	})
	require.Error(t, err)
	assert.EqualError(t, err, "name cannot be empty")
	assert.Equal(t, 0, id)
	repo.AssertNotCalled(t, "Add")
}
func TestAdminProductUseCase_Add_InvalidPrice(t *testing.T) {

	repo := new(MockProductAdminRepository)
	uc := AdminProductUseCase{productRepo: repo}

	id, err := uc.Add(domain.Product{
		Name:  "Ноутбук",
		Price: 0,
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrPrice)
	assert.Equal(t, 0, id)
	repo.AssertNotCalled(t, "Add")
}
func TestAdminProductUseCase_Add_RepoError(t *testing.T) {

	repo := new(MockProductAdminRepository)
	uc := AdminProductUseCase{productRepo: repo}

	repo.On("Add", mock.Anything).Return(0, errors.New("ошибка базы данных"))

	id, err := uc.Add(domain.Product{
		Name:  "Ноутбук",
		Price: 1000,
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка базы данных")
	assert.Equal(t, 0, id)
	repo.AssertExpectations(t)
}
func TestUpdate_Success(t *testing.T) {

	adminRepo := new(MockProductAdminRepository)
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		productRepo: adminRepo,
		catalogRepo: catalogRepo,
	}

	testProduct := domain.Product{
		ID:    1,
		Name:  "Updated Product",
		Price: 2000,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: false,
	}, nil)

	adminRepo.On("Update", testProduct).Return(nil)

	err := uc.Update(testProduct)

	assert.NoError(t, err)
	adminRepo.AssertExpectations(t)
	catalogRepo.AssertExpectations(t)
}

func TestUpdate_InvalidID(t *testing.T) {
	uc := AdminProductUseCase{}

	err := uc.Update(domain.Product{
		ID:    0,
		Name:  "Product",
		Price: 1000,
	})
	assert.ErrorIs(t, err, domain.ErrInvalidIdProduct)
}

func TestUpdate_EmptyName(t *testing.T) {
	uc := AdminProductUseCase{}

	err := uc.Update(domain.Product{
		ID:    1,
		Name:  "",
		Price: 1000,
	})
	assert.ErrorIs(t, err, domain.ErrName)
}
func TestUpdate_InvalidPrice(t *testing.T) {
	uc := AdminProductUseCase{}

	err := uc.Update(domain.Product{
		ID:    1,
		Name:  "Product",
		Price: 0,
	})
	assert.ErrorIs(t, err, domain.ErrPrice)
}
func TestUpdate_ProductNotFound(t *testing.T) {
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{}, domain.ErrProductNotFound)

	err := uc.Update(domain.Product{
		ID:    1,
		Name:  "Product",
		Price: 1000,
	})
	assert.ErrorContains(t, err, "product not found")
	assert.ErrorIs(t, err, domain.ErrProductNotFound)
}

func TestUpdate_ArchivedProduct(t *testing.T) {
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: true,
	}, nil)

	err := uc.Update(domain.Product{
		ID:    1,
		Name:  "Product",
		Price: 1000,
	})

	assert.ErrorIs(t, err, domain.ErrUpdateArchived)
}
func TestUpdate_RepositoryError(t *testing.T) {
	adminRepo := new(MockProductAdminRepository)
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		productRepo: adminRepo,
		catalogRepo: catalogRepo,
	}

	testProduct := domain.Product{
		ID:    1,
		Name:  "Product",
		Price: 1000,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: false,
	}, nil)

	adminRepo.On("Update", testProduct).Return(errors.New("update failed"))

	err := uc.Update(testProduct)

	assert.ErrorContains(t, err, "update failed")
}
func TestArchiv_Success(t *testing.T) {

	adminRepo := new(MockProductAdminRepository)
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		productRepo: adminRepo,
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: false,
	}, nil)

	adminRepo.On("Archiv", 1).Return(nil)

	err := uc.Archiv(1)

	assert.NoError(t, err)
	adminRepo.AssertExpectations(t)
	catalogRepo.AssertExpectations(t)
}

func TestArchiv_InvalidID(t *testing.T) {
	uc := AdminProductUseCase{}
	err := uc.Archiv(0)
	assert.ErrorIs(t, err, domain.ErrInvalidIdProduct)
}
func TestArchiv_ProductNotFound(t *testing.T) {
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{}, domain.ErrProductNotFound)

	err := uc.Archiv(1)

	assert.ErrorIs(t, err, domain.ErrInvalidIdProduct)
	assert.EqualError(t, err, "ID must be positive")
}

func TestArchiv_AlreadyArchived(t *testing.T) {
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: true,
	}, nil)

	err := uc.Archiv(1)

	assert.ErrorIs(t, err, domain.ErrProductArchived)
}
func TestArchiv_RepositoryError(t *testing.T) {
	adminRepo := new(MockProductAdminRepository)
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		productRepo: adminRepo,
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: false,
	}, nil)

	adminRepo.On("Archiv", 1).Return(errors.New("database error"))

	err := uc.Archiv(1)

	assert.ErrorContains(t, err, "database error")
}

func TestRestore_Success(t *testing.T) {
	adminRepo := new(MockProductAdminRepository)
	catalogRepo := new(MockCatalogRepository)

	uc := AdminProductUseCase{
		productRepo: adminRepo,
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: true,
	}, nil)

	adminRepo.On("Restore", 1).Return(nil)

	err := uc.Restore(1)

	require.NoError(t, err)
	adminRepo.AssertExpectations(t)
	catalogRepo.AssertExpectations(t)
}
func TestRestore_InvalidID(t *testing.T) {
	uc := AdminProductUseCase{}

	err := uc.Restore(0)

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidIdProduct)
	assert.EqualError(t, err, "ID must be positive")
}

func TestRestore_ProductNotFound(t *testing.T) {
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{}, domain.ErrProductNotFound)

	err := uc.Restore(1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidIdProduct)
	assert.Contains(t, err.Error(), "ID must be positive")
}
func TestRestore_NotArchived(t *testing.T) {
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: false,
	}, nil)

	err := uc.Restore(1)

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrProductArchived)
	assert.EqualError(t, err, "product on archive")
}
func TestRestore_RepositoryError(t *testing.T) {
	adminRepo := new(MockProductAdminRepository)
	catalogRepo := new(MockCatalogRepository)
	uc := AdminProductUseCase{
		productRepo: adminRepo,
		catalogRepo: catalogRepo,
	}

	catalogRepo.On("GetByID", 1).Return(domain.Product{
		ID:         1,
		IsArchived: true,
	}, nil)

	adminRepo.On("Restore", 1).Return(errors.New("database error"))

	err := uc.Restore(1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
}

func TestGetArchived_Success(t *testing.T) {
	repo := new(MockProductAdminRepository)
	uc := AdminProductUseCase{productRepo: repo}

	expectedProducts := []domain.Product{
		{ID: 1, Name: "Archived Product 1", IsArchived: true},
		{ID: 2, Name: "Archived Product 2", IsArchived: true},
	}
	repo.On("GetArchived").Return(expectedProducts, nil)

	products, err := uc.GetArchived()

	require.NoError(t, err)
	assert.Equal(t, expectedProducts, products)
	repo.AssertExpectations(t)
}
func TestGetArchived_EmptyList(t *testing.T) {
	repo := new(MockProductAdminRepository)
	uc := AdminProductUseCase{productRepo: repo}

	repo.On("GetArchived").Return([]domain.Product{}, nil)

	products, err := uc.GetArchived()

	require.NoError(t, err)
	assert.Empty(t, products)
	repo.AssertExpectations(t)
}

func TestGetArchived_RepositoryError(t *testing.T) {
	repo := new(MockProductAdminRepository)
	uc := AdminProductUseCase{productRepo: repo}

	expectedErr := errors.New("database connection failed")
	repo.On("GetArchived").Return(nil, expectedErr)

	products, err := uc.GetArchived()

	require.Error(t, err)
	assert.Nil(t, products)
	assert.Contains(t, err.Error(), "error getting archive product")
	assert.ErrorIs(t, err, expectedErr)
	repo.AssertExpectations(t)
}



package postgres_test

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
	"github.com/Vladislav5108/inventorytracker/internal/infrastructure/db/postgres"
)

func newMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return db, mock
}
func TestGetByID_Success(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	testID := 1

	expectedProduct := domain.Product{
		ID:         testID,
		Name:       "MacBook Pro",
		Price:      199999,
		Quantity:   10,
		CategoryID: 3,
		CreatedAt:  fixedTime,
		IsArchived: false,
	}
	expectedQuery := `SELECT id, name, price, quantity, category_id, created_at, is_archived
	FROM products
	WHERE id = \$1`

	mock.ExpectQuery(expectedQuery).
		WithArgs(testID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "price", "quantity", "category_id", "created_at", "is_archived",
		}).AddRow(
			expectedProduct.ID,
			expectedProduct.Name,
			expectedProduct.Price,
			expectedProduct.Quantity,
			expectedProduct.CategoryID,
			expectedProduct.CreatedAt,
			expectedProduct.IsArchived,
		))
	product, err := repo.GetByID(testID)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(product, expectedProduct) {
		t.Errorf("expected %+v, got %+v", expectedProduct, product)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetByID_NotFound(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	nonExistentID := 999

	expectedQuery := `SELECT id, name, price, quantity, category_id, created_at, is_archived
	FROM products
	WHERE id = \$1`

	mock.ExpectQuery(expectedQuery).
		WithArgs(nonExistentID).
		WillReturnError(sql.ErrNoRows)

	_, err := repo.GetByID(nonExistentID)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("expected ErrProductNotFound, got %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}
func TestGetByID_DatabaseError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testID := 1
	expectedErr := fmt.Errorf("conntction timeout")

	expectedQuery := `SELECT id, name, price, quantity, category_id, created_at, is_archived
	FROM products
	WHERE id = \$1`

	mock.ExpectQuery(expectedQuery).
		WithArgs(testID).
		WillReturnError(expectedErr)

	_, err := repo.GetByID(testID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), fmt.Sprintf("failed to get product ID %d", testID)) {
		t.Errorf("error should contain product ID, got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected underlying error %v, got %v", expectedErr, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %s", err)
	}
}

func TestGetAll_Success(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	expectedProducts := []domain.Product{
		{
			ID:         1,
			Name:       "MacBook Pro",
			Price:      199999,
			Quantity:   10,
			CategoryID: 3,
			CreatedAt:  fixedTime,
			IsArchived: false,
		},
		{
			ID:         2,
			Name:       "iPhone 15",
			Price:      899999,
			Quantity:   5,
			CategoryID: 4,
			CreatedAt:  fixedTime,
			IsArchived: false,
		},
	}
	query := regexp.QuoteMeta(`
	SELECT id, name, price, quantity, category_id, 
	created_at, is_archived FROM products`)

	mock.ExpectQuery(query).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "price", "quantity",
			"category_id", "created_at", "is_archived",
		}).AddRow(
			expectedProducts[0].ID,
			expectedProducts[0].Name,
			expectedProducts[0].Price,
			expectedProducts[0].Quantity,
			expectedProducts[0].CategoryID,
			expectedProducts[0].CreatedAt,
			expectedProducts[0].IsArchived,
		).AddRow(
			expectedProducts[1].ID,
			expectedProducts[1].Name,
			expectedProducts[1].Price,
			expectedProducts[1].Quantity,
			expectedProducts[1].CategoryID,
			expectedProducts[1].CreatedAt,
			expectedProducts[1].IsArchived,
		))
	products, err := repo.GetAll()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(products, expectedProducts) {
		t.Errorf("expected %v, got %v", expectedProducts, products)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAll_DatabaseError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	expectedErr := fmt.Errorf("database connection failed")
	query := regexp.QuoteMeta(`
        SELECT id, name, price, quantity, category_id, 
        created_at, is_archived FROM products`)

	mock.ExpectQuery(query).
		WillReturnError(expectedErr)

	products, err := repo.GetAll()

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "list request error") {
		t.Errorf("error should contain 'list request error', got: %v", err)
	}
	if products != nil {
		t.Errorf("expected nil products, got %v", products)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestGetByCategory_Success(t *testing.T) {

	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	categoryID := 3
	testTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	expectedProducts := []domain.Product{
		{
			ID:         1,
			Name:       "MacBook Pro",
			Price:      199999,
			Quantity:   10,
			CategoryID: categoryID,
			CreatedAt:  testTime,
			IsArchived: false,
		},
	}
	expectedSQL := `SELECT id, name, price, quantity, category_id, created_at, is_archived
FROM products
WHERE category_id = \$1`

	mock.ExpectQuery(expectedSQL).
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "price", "quantity", "category_id", "created_at", "is_archived",
		}).AddRow(
			expectedProducts[0].ID,
			expectedProducts[0].Name,
			expectedProducts[0].Price,
			expectedProducts[0].Quantity,
			expectedProducts[0].CategoryID,
			expectedProducts[0].CreatedAt,
			expectedProducts[0].IsArchived,
		))

	products, err := repo.GetByCategory(categoryID)

	if err != nil {
		t.Fatalf("Unexpected error occurred: %v", err)
	}
	if len(products) != len(expectedProducts) {
		t.Fatalf("Expected %d items, got %d",
			len(expectedProducts), len(products))
	}
	if !reflect.DeepEqual(products, expectedProducts) {
		t.Errorf("Expected: %+v\ngot: %+v", expectedProducts, products)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}
func TestGetByCategory_Empty(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	categoryID := 999
	mock.ExpectQuery(`SELECT id, name, price, quantity, category_id, created_at, is_archived
FROM products
WHERE category_id = \$1`).
		WithArgs(categoryID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "price", "quantity", "category_id", "created_at", "is_archived",
		}))

	products, err := repo.GetByCategory(categoryID)

	if err != nil {
		t.Errorf("Expected nil error, got: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("Expected empty list, got %d products", len(products))
	}
}
func TestGetByCategory_DBError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	categoryID := 3

	expectedErr := fmt.Errorf("connection timeout")

	mock.ExpectQuery(`SELECT id, name, price, quantity, category_id, created_at, is_archived
FROM products
WHERE category_id = \$1`).
		WithArgs(categoryID).
		WillReturnError(expectedErr)

	products, err := repo.GetByCategory(categoryID)

	if err == nil {
		t.Fatalf("Expected an error, but got nil")
	}
	if !strings.Contains(err.Error(), "category list request error") {
		t.Errorf("Expected error containing 'category list request error', got: %v", err)
	}
	if products != nil {
		t.Error("Expected nil products when error occurs")
	}
}
func TestAdd_Success(t *testing.T) {

	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testProduct := domain.Product{
		Name:       "MacBook Pro",
		Price:      199999,
		Quantity:   10,
		CategoryID: 3,
		IsArchived: false,
	}
	expectedID := 123

	expectedSQL := `INSERT INTO products\(
	name, price, quantity, category_id, is_archived, created_at\)
	VALUES\(\$1,\$2,\$3,\$4,\$5,\$6\)
	RETURNING id`

	mock.ExpectQuery(expectedSQL).
		WithArgs(
			testProduct.Name,
			testProduct.Price,
			testProduct.Quantity,
			testProduct.CategoryID,
			testProduct.IsArchived,
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	id, err := repo.Add(testProduct)

	if err != nil {
		t.Fatalf("No error was expected, but got: %v", err)
	}
	if id != expectedID {
		t.Errorf("Expected ID = %d, got %d", expectedID, id)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}
func TestAdd_DBError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testProduct := domain.Product{
		Name:       "iPhone 13",
		Price:      89999,
		Quantity:   5,
		CategoryID: 4,
		IsArchived: false,
	}
	expectedErr := fmt.Errorf("connection failed")

	mock.ExpectQuery(`INSERT INTO products\(
	name, price, quantity, category_id, is_archived, created_at\)
	VALUES\(\$1,\$2,\$3,\$4,\$5,\$6\)
	RETURNING id`).
		WithArgs(
			testProduct.Name,
			testProduct.Price,
			testProduct.Quantity,
			testProduct.CategoryID,
			testProduct.IsArchived,
			sqlmock.AnyArg(),
		).
		WillReturnError(expectedErr)

	id, err := repo.Add(testProduct)

	if err == nil {
		t.Fatal("xpected an error, but got nil")
	}
	if !strings.Contains(err.Error(), "error add data") {
		t.Errorf("Expected error to contain 'error add data', got: %v", err)
	}
	if id != 0 {
		t.Errorf("Expected ID = 0 on error, got %d", id)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}
func TestAdd_InvalidData(t *testing.T) {

	db, _ := newMock(t)
	repo := postgres.NewProductRepo(db)

	testCases := []struct {
		name    string
		product domain.Product
		errMsg  string
	}{
		{
			name: "Empty name",
			product: domain.Product{
				Name:       "",
				Price:      1000,
				Quantity:   1,
				CategoryID: 1,
			},
			errMsg: "name",
		},
		{
			name: "Negative price",
			product: domain.Product{
				Name:       "Product",
				Price:      -100,
				Quantity:   1,
				CategoryID: 1,
			},
			errMsg: "price",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := repo.Add(tc.product)

			if err == nil {
				t.Fatal("expected an error, but got nil")
			}
			if !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("Expected error containing '%s', got: %v",
					tc.errMsg, err)
			}
		})
	}
}

func TestUpdate_Success(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testProduct := domain.Product{
		ID:         1,
		Name:       "MacBook Pro Update",
		Price:      2099999,
		Quantity:   15,
		CategoryID: 3,
	}
	expectedSQL := `UPDATE products SET
	name = \$1, price = \$2, quantity = \$3, category_id = \$4
	WHERE id = \$5`

	mock.ExpectExec(expectedSQL).
		WithArgs(
			testProduct.Name,
			testProduct.Price,
			testProduct.Quantity,
			testProduct.CategoryID,
			testProduct.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Update(testProduct)

	if err != nil {
		t.Fatalf("No error was expected, but got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}
func TestUpdate_NotFound(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testProduct := domain.Product{
		ID:         999,
		Name:       "Not Found Product",
		Price:      1000,
		Quantity:   1,
		CategoryID: 1,
	}

	mock.ExpectExec(`UPDATE products SET
	name = \$1, price = \$2, quantity = \$3, category_id = \$4
	WHERE id = \$5`).
		WithArgs(
			testProduct.Name,
			testProduct.Price,
			testProduct.Quantity,
			testProduct.CategoryID,
			testProduct.ID,
		).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Update(testProduct)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("Expected error ErrProductNotFound, got: %v", err)
	}
}
func TestUpdate_DBError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testProduct := domain.Product{
		ID:         1,
		Name:       "MacBook Pro",
		Price:      199999,
		Quantity:   10,
		CategoryID: 3,
	}
	expectedErr := fmt.Errorf("connection timeout")

	mock.ExpectExec(`UPDATE products SET
	name = \$1, price = \$2, quantity = \$3, category_id = \$4
	WHERE id = \$5`).
		WithArgs(
			testProduct.Name,
			testProduct.Price,
			testProduct.Quantity,
			testProduct.CategoryID,
			testProduct.ID,
		).
		WillReturnError(expectedErr)

	err := repo.Update(testProduct)

	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}

	if !strings.Contains(err.Error(), "error update data") {
		t.Errorf("Expected error to contain 'error update data', got: %v", err)
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestArchiv_Success(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testID := 1

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE products SET is_archived = true WHERE id = $1 AND is_archived = false`)).
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Archiv(testID)

	if err != nil {
		t.Fatalf("No error was expected, but one occurred: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}

func TestArchiv_NotFound(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testID := 999

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE products SET is_archived = true WHERE id = $1 AND is_archived = false`)).
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Archiv(testID)
	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("Expected error ErrProductNotFound, got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}
func TestArchiv_DBError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testID := 1
	expectedErr := fmt.Errorf("connection failed")

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE products SET is_archived = true WHERE id = $1 AND is_archived = false`)).
		WithArgs(testID).
		WillReturnError(expectedErr)

	err := repo.Archiv(testID)

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	if !strings.Contains(err.Error(), "archive error") {
		t.Errorf("Expected error on text 'archive error', got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestRestore_Success(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testID := 1

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE products SET is_archived = false WHERE id = $1`)).
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := repo.Restore(testID)

	if err != nil {
		t.Fatalf("Error was not expected but received")
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectation not met: %v", err)
	}
}
func TestRestore_NotFound(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testID := 999

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE products SET is_archived = false WHERE id = $1`)).
		WithArgs(testID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.Restore(testID)

	if !errors.Is(err, domain.ErrProductNotFound) {
		t.Errorf("Expected error ErrProductNotFound, got: %v", err)
	}
}
func TestRestore_DBError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	testID := 1
	expectedErr := fmt.Errorf("connection failed")

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE products SET is_archived = false WHERE id = $1`)).
		WithArgs(testID).
		WillReturnError(expectedErr)

	err := repo.Restore(testID)

	if err == nil {
		t.Fatal("Expected on error, but got nil")
	}
	if !strings.Contains(err.Error(), "error delete on archive") {
		t.Errorf("Expected error with text 'error delete on archive', got: %v", err)
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected on error %v, got %v", expectedErr, err)
	}
}
func TestGetArchived_Success(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	expectedProducts := []domain.Product{
		{
			ID:         1,
			Name:       "Archived Product 1",
			Price:      10000,
			Quantity:   5,
			CategoryID: 2,
			CreatedAt:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			IsArchived: true,
		},
		{
			ID:         2,
			Name:       "Archived Product 2",
			Price:      20000,
			Quantity:   3,
			CategoryID: 3,
			CreatedAt:  time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
			IsArchived: true,
		},
	}

	expectedQuery := `SELECT id, name, price, quantity, category_id, created_at, is_archived 
    FROM products
    WHERE is_archived = true`

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "price", "quantity", "category_id", "created_at", "is_archived",
		}).AddRow(
			expectedProducts[0].ID,
			expectedProducts[0].Name,
			expectedProducts[0].Price,
			expectedProducts[0].Quantity,
			expectedProducts[0].CategoryID,
			expectedProducts[0].CreatedAt,
			expectedProducts[0].IsArchived,
		).AddRow(
			expectedProducts[1].ID,
			expectedProducts[1].Name,
			expectedProducts[1].Price,
			expectedProducts[1].Quantity,
			expectedProducts[1].CategoryID,
			expectedProducts[1].CreatedAt,
			expectedProducts[1].IsArchived,
		))

	products, err := repo.GetArchived()

	if err != nil {
		t.Fatalf("Error was not expected but received: %v", err)
	}
	if len(products) != len(expectedProducts) {
		t.Fatalf("Expected %d products, got %d", len(expectedProducts), len(products))
	}
	if !reflect.DeepEqual(products, expectedProducts) {
		t.Errorf("Expected: %+v\nGot: %+v", expectedProducts, products)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unmet mock expectations: %v", err)
	}
}
func TestGetArchived_Empty(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, price, quantity, category_id, created_at, is_archived 
	FROM products
	WHERE is_archived = true`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "price", "quantity", "category_id", "created_at", "is_archived",
		}))
	products, err := repo.GetArchived()

	if err != nil {
		t.Fatalf("Unexpected error received: %v", err)
	}
	if len(products) != 0 {
		t.Errorf("Expected empty list, got %d items", len(products))
	}
}
func TestGetArchived_DBError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	expectedErr := fmt.Errorf("database connection error")

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, price, quantity, category_id, created_at, is_archived 
	FROM products
	WHERE is_archived = true`)).
		WillReturnError(expectedErr)

	products, err := repo.GetArchived()

	if err == nil {
		t.Fatal("Expected an error but got nil")
	}
	if !strings.Contains(err.Error(), "error getting the archive list") {
		t.Errorf("Expected error with text 'error getting the archive list', got: %v", err)
	}
	if products != nil {
		t.Error("Expected nil products on error")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error: %v, got: %v", expectedErr, err)
	}
}
func TestGetArchived_ScanError(t *testing.T) {
	db, mock := newMock(t)
	repo := postgres.NewProductRepo(db)

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, price, quantity, category_id, created_at, is_archived 
	FROM products
	WHERE is_archived = true`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "name", "price", "quantity", "category_id", "created_at", "is_archived",
		}).AddRow(
			1, "Product", "invalid price", 10, 2, time.Now(), true,
		))

	_, err := repo.GetArchived()

	if err == nil {
		t.Fatal("Expected an error but got nil")
	}
	if !strings.Contains(err.Error(), "error read data") {
		t.Errorf("Expected an error with text'error read data', got: %v", err)
	}
}

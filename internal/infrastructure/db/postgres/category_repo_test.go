package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	domain "github.com/Vladislav5108/inventorytracker/internal/domain/entity"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}
	return db, mock
}

func TestCreateCategory(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := NewCategoryRepo(db)

	mock.ExpectQuery(`INSERT INTO categories`).
		WithArgs("Ноутбуки", "Техника для работы").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	id, err := repo.CreateCategory(domain.Category{
		Name:        "Ноутбуки",
		Description: "Техника для работы",
	})

	if err != nil {
		t.Errorf("No error was expected, but we got one: %v", err)
	}
	if id != 123 {
		t.Errorf("Expected id=123, got something else %d", id)
	}
}

func TestCreateCategory_DBError(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := NewCategoryRepo(db)

	mock.ExpectQuery(`INSERT INTO caregory`).
		WithArgs("Телефоны", "").
		WillReturnError(fmt.Errorf("error connection"))

	_, err := repo.CreateCategory(domain.Category{
		Name:        "Телефоны",
		Description: "",
	})

	if err == nil {
		t.Error("We expected an error, but there is none!")
	}
}
func TestGetByIDCategory_Success(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	repo := NewCategoryRepo(db)

	expectedQuery := regexp.QuoteMeta(`SELECT id, name, description FROM categories WHERE id = $1`)

	mock.ExpectQuery(expectedQuery).
		WithArgs(1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, "Ноутбуки", "Техника для работы"),
		)

	category, err := repo.GetByIDCategory(1)

	if err != nil {
		t.Fatalf("No error was expected, but we got one: %v", err)
	}

	expected := domain.Category{
		ID:          1,
		Name:        "Ноутбуки",
		Description: "Техника для работы",
	}
	if !reflect.DeepEqual(category, expected) {
		t.Fatalf("Expected %+v, got %+v", expected, category)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestGetByIDCategory_NotFound(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	expectedQuery := `SELECT id, name, description FROM categories WHERE id = \$1`

	mock.ExpectQuery(expectedQuery).
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	repo := NewCategoryRepo(db)

	result, err := repo.GetByIDCategory(999)
	if !errors.Is(err, domain.ErrNotFoundCategory) {
		t.Errorf("Expected error ErrNotFoundCategory, got: %v", err)
	}
	if (result != domain.Category{}) {
		t.Errorf("Expected an empty category, got one: %v", result)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations were not met: %v", err)
	}
}

func TestGetAllCategories_Success(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	testCategories := []domain.Category{
		{ID: 1, Name: "Ноутбуки", Description: "Техника для работы"},
		{ID: 2, Name: "Смартфоны", Description: "Мобильные устройства"},
	}

	expectedQuery := `SELECT id, name, description FROM categories`

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "Ноутбуки", "Техника для работы").
		AddRow(2, "Смартфоны", "Мобильные устройства")

	mock.ExpectQuery(expectedQuery).WillReturnRows(rows)

	repo := NewCategoryRepo(db)
	result, err := repo.GetAllCategories()

	if err != nil {
		t.Fatalf("No error was expected, but we got one: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("Got 2 categories as expected %d", len(result))
	}

	if !reflect.DeepEqual(result, testCategories) {
		t.Fatalf("Expected %+v, got %+v", testCategories, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Database expectations were not met: %v", err)
	}
}

func TestGetAllCategories_DBError(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(`SELECT id, name, description FROM categories`).
		WillReturnError(fmt.Errorf("error database"))

	repo := NewCategoryRepo(db)

	result, err := repo.GetAllCategories()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
	if result != nil {
		t.Fatalf("Expected nil, got %v", result)
	}
}
func TestGetAllCategories_ScanError(t *testing.T) {

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectQuery(`SELECT id, name, description FROM categories`).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(1, "Ноутбуки", nil),
		)

	repo := NewCategoryRepo(db)

	_, err := repo.GetAllCategories()

	if err == nil {
		t.Fatal("Expected error, but got nil")
	}
}
func TestUpdateCategory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error: %v", err)
	}
	defer db.Close()

	repo := &CategoryRepo{db: db}

	testCategory := domain.Category{
		ID:          1,
		Name:        "Обновленное имя",
		Description: "Обновленное описание",
	}
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, description FROM categories WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(1, "Старое имя", "Старое описание"))

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE categories SET name = $1, description = $2 WHERE id = $3`)).
		WithArgs(testCategory.Name, testCategory.Description, testCategory.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateCategory(testCategory)
	if err != nil {
		t.Errorf("No error was expected, but we got one: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestUpdateCategory_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error: %v", err)
	}
	defer db.Close()

	repo := CategoryRepo{db: db}

	testCategory := domain.Category{ID: 999}
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, description FROM categories WHERE id = $1`)).
		WithArgs(999).
		WillReturnError(domain.ErrNotFoundCategory)

	err = repo.UpdateCategory(testCategory)
	if !errors.Is(err, domain.ErrNotFoundCategory) {
		t.Errorf("Expected error ErrNotFoundCategory, got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestUpdateCategory_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error:%v", err)
	}
	defer db.Close()

	repo := &CategoryRepo{db: db}

	testCategory := domain.Category{
		ID:          1,
		Name:        "Новое имя",
		Description: "Новое описание",
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, description FROM categories WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(1, "Старое имя", "Старое описание"))

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE categories SET name = $1, description = $2 WHERE id = $3`)).
		WithArgs(testCategory.Name, testCategory.Description, testCategory.ID).
		WillReturnError(fmt.Errorf("database error"))

	err = repo.UpdateCategory(testCategory)

	if err == nil {
		t.Error("Expected error, but got nil")
	}
	if !strings.Contains(err.Error(), "error update category") {
		t.Errorf("Expected an update error, but it completed successfully: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}

func TestUpdateCategory_NoRowsAffected(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error: %v", err)
	}
	defer db.Close()

	repo := &CategoryRepo{db: db}

	testCategory := domain.Category{
		ID:          1,
		Name:        "Новое имя",
		Description: "Новое описание",
	}
	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, description FROM categories WHERE id = $1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(1, "Старое имя", "Старое описание"))

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE categories SET name = $1, description = $2 WHERE id = $3`)).
		WithArgs(testCategory.Name, testCategory.Description, testCategory.ID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.UpdateCategory(testCategory)

	if !errors.Is(err, domain.ErrNotFoundCategory) {
		t.Errorf("Expected error ErrNotFoundCategory, got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Database expectations were not met: %v", err)
	}
}
func TestDeleteCategory_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error: %v", err)
	}
	defer db.Close()

	repo := &CategoryRepo{db: db}

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM categories WHERE id = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteCategory(1)

	if err != nil {
		t.Errorf("No error was expected, but we got one.: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %v", err)
	}
}
func TestDeleteCategory_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error: %v", err)
	}
	defer db.Close()

	repo := &CategoryRepo{db: db}

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM categories WHERE id = $1`)).
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = repo.DeleteCategory(999)

	if !errors.Is(err, domain.ErrNotFoundCategory) {
		t.Errorf("Expected ErrNotFoundCategory, got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %v", err)
	}
}

func TestDeleteCategory_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error: %v", err)
	}
	defer db.Close()

	repo := &CategoryRepo{db: db}

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM categories WHERE id = $1`)).
		WithArgs(1).
		WillReturnError(fmt.Errorf("error connect database"))

	err = repo.DeleteCategory(1)

	if err == nil {
		t.Error("expected error, but got nil")
	}
	if !strings.Contains(err.Error(), "error delete") {
		t.Errorf("expected error 'error delete', got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %v", err)
	}
}

func TestDeleteCategory_RowsAffectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Mock creation error: %v", err)
	}
	defer db.Close()

	repo := &CategoryRepo{db: db}

	mock.ExpectExec(regexp.QuoteMeta(
		`DELETE FROM categories WHERE id = $1`)).
		WithArgs(1).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("error RowsAffected")))

	err = repo.DeleteCategory(1)

	if err == nil {
		t.Error("Expected an error but got nil")
	}
	if !strings.Contains(err.Error(), "error delete") {
		t.Errorf("expected text 'error delete', got: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %v", err)
	}
}

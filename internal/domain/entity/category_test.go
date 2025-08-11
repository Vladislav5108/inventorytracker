package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategory_Validate_Success(t *testing.T) {
	category := Category{
		Name: "смартфоны",
	}
	err := category.Validate()

	assert.NoError(t, err, "Validation was expected to pass without errors")
}
func Test_Category_Validate_EmptyName(t *testing.T) {
	category := Category{
		Name: "",
	}
	err := category.Validate()
	assert.ErrorIs(t, err, ErrNameCategory, "Expected on error 'ErrNameCategory'")
	assert.Equal(t, "category name cannot be empty", err.Error(), "incorrect error message")

}

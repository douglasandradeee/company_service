package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testes para ValidCNPJ
func TestValidCNPJ_Empty_ReturnsFalse(t *testing.T) {
	cnpj := ""
	assert.False(t, ValidCNPJ(cnpj))
}

func TestValidCNPJ_TooShort_ReturnsFalse(t *testing.T) {
	cnpj := "123"
	assert.False(t, ValidCNPJ(cnpj))
}

func TestValidCNPJ_TooLong_ReturnsFalse(t *testing.T) {
	cnpj := "123456789012345"
	assert.False(t, ValidCNPJ(cnpj))
}

func TestValidCNPJ_AllSameDigits_ReturnsFalse(t *testing.T) {
	assert.False(t, ValidCNPJ("00000000000000"))
	assert.False(t, ValidCNPJ("11111111111111"))
}

func TestValidCNPJ_ValidFormat_ReturnsTrue(t *testing.T) {
	assert.True(t, ValidCNPJ("11444777000161"))
	assert.True(t, ValidCNPJ("12345678901234")) // Apenas verifica formato, não dígitos
}

func TestCleanCNPJ_RemovesNonDigits(t *testing.T) {
	assert.Equal(t, "11444777000161", CleanCNPJ("11.444.777/0001-61"))
	assert.Equal(t, "11444777000161", CleanCNPJ("11-444-777-0001-61"))
	assert.Equal(t, "11444777000161", CleanCNPJ("11 444 777 0001 61"))
}

func TestCleanCNPJ_EmptyString_ReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", CleanCNPJ(""))
}

func TestCleanCNPJ_OnlyNonDigits_ReturnsEmpty(t *testing.T) {
	assert.Equal(t, "", CleanCNPJ("abc.-/"))
}

func TestAllDigitsEqual_Empty_ReturnsFalse(t *testing.T) {
	assert.False(t, hasDifferentDigits(""))
}

func TestFormatCNPJ_Valid_ReturnsFormatted(t *testing.T) {
	assert.Equal(t, "11.444.777/0001-61", FormatCNPJ("11444777000161"))
}

func TestFormatCNPJ_TooShort_ReturnsOriginal(t *testing.T) {
	assert.Equal(t, "123", FormatCNPJ("123"))
}

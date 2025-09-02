package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGivenEmptyCNPJ_WhenInValidCNPJ_ThenReturnFalse(t *testing.T) {
	cnpj := ""
	assert.False(t, ValidCNPJ(cnpj))
}

func TestGivenCNPJ_WhenInvalidLength_ThenReturnFalse(t *testing.T) {
	cnpj := "12345678901234"
	assert.False(t, ValidCNPJ(cnpj))
}

func TestGivenCNPJ_WhenAllDigitsEqual_ThenReturnFalse(t *testing.T) {
	cnpj := "11111111111111"
	assert.False(t, ValidCNPJ(cnpj))
}

func TestGivenCNPJ_WhenInvalidCheckDigits_ThenReturnFalse(t *testing.T) {
	cnpj := "12345678901234"
	assert.False(t, ValidCNPJ(cnpj))
}

package domain

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGivenCompany_WhenInvalidCNPJ_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
	}

	// Then
	assert.Error(t, company.Validate(), "Expected error for invalid CNPJ")
}

func TestGivenCompany_WhenInvalidFantasyName_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
	}

	// Then
	assert.Error(t, company.Validate(), "Expected error for invalid Fantasy Name")
}

func TestGivenCompany_WhenFantasyNameLenHasMinLimit_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "A",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
	}

	// Then
	assert.Error(t, company.validateNames(), "Expected no error for valid Fantasy Name")
}

func TestGivenCompany_WhenFantasyNameLenHasMaxLimit_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 strings.Repeat("a", 101),
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
	}

	// Then
	assert.Error(t, company.validateNames(), "Expected no error for valid Fantasy Name")
}

func TestGivenCompany_WhenInvalidCorporateName_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
	}

	// Then
	assert.Error(t, company.Validate(), "Expected error for invalid Corporate Name")
}

func TestGivenCompany_WhenCorporateNameLenHasMinLimit_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste LTDA",
		CorporateName:               "A",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
	}

	// Then
	assert.Error(t, company.validateNames(), "Expected no error for valid Fantasy Name")
}

func TestGivenCompany_WhenCorporateNameLenHasMaxLimit_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste LTDA",
		CorporateName:               strings.Repeat("a", 151),
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
	}

	// Then
	assert.Error(t, company.validateNames(), "Expected no error for valid Fantasy Name")
}

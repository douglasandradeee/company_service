package domain

import (
	"strings"
	"testing"
	"time"

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
	assert.Error(t, company.validateNames(), "Expected error for invalid Fantasy Name")
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
	assert.Error(t, company.validateNames(), "Expected error for invalid Fantasy Name")
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
	assert.Error(t, company.validateNames(), "Expected error for invalid Corporate Name")
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
	assert.Error(t, company.validateNames(), "Expected error for invalid Corporate Name")
}

func TestGivenCompany_WhenLargeCompanyWithoutPCD_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               100, // 100+ funcionários
		RequiredMinPWDEmployeeCount: 0,   // Sem PCD
	}

	// When
	err := company.Validate()

	// Then
	assert.Error(t, err, "Expected error for invalid PCD count")
}

func TestCompany_BeforeCreate_SetsTimestamps(t *testing.T) {
	// Given
	company := &Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
	}

	// When
	company.BeforeCreate()

	// Then
	assert.False(t, company.CreatedAt.IsZero(), "CreatedAt should be set")
	assert.False(t, company.UpdatedAt.IsZero(), "UpdatedAt should be set")
	assert.Equal(t, company.CreatedAt, company.UpdatedAt, "Timestamps should be equal on creation")
}

func TestCompany_BeforeUpdate_UpdatesOnlyUpdatedAt(t *testing.T) {
	// Given
	createdAt := time.Now().Add(-time.Hour * 24)
	company := &Company{
		ID:                          "1",
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 1,
		CreatedAt:                   createdAt,
		UpdatedAt:                   createdAt,
	}

	// When
	company.BeforeUpdate()

	// Then
	assert.True(t, company.UpdatedAt.After(createdAt), "UpdatedAt should be updated")
	assert.Equal(t, createdAt, company.CreatedAt, "CreatedAt should not be modified")
}

func TestGivenCompany_WhenCountEmployeesIsNegative_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               -1,
		RequiredMinPWDEmployeeCount: 1,
	}

	// When
	err := company.Validate()

	// Then
	assert.Error(t, err, "Expected error for invalid Employee Count")
}

func TestGivenCompany_WhenCountEmployeesIsZero_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               0,
		RequiredMinPWDEmployeeCount: 1,
	}

	// When
	err := company.Validate()

	// Then
	assert.Error(t, err, "Expected error for invalid Employee Count")
}

func TestGivenCompany_WhenPWDEmployeeIsGreaterThanEmployeeCount_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123, Centro - São Paulo/SP",
		EmployeeCount:               5,
		RequiredMinPWDEmployeeCount: 10,
	}

	// When
	err := company.Validate()

	// Then
	assert.Error(t, err, "Expected error for PWD Employee Count greater than Employee Count")
}

func TestGivenCompany_WhenAddressIsInvalid_ThenShouldReturnError(t *testing.T) {
	// Given
	company := Company{
		CNPJ:                        "11444777000161",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "",
		EmployeeCount:               5,
		RequiredMinPWDEmployeeCount: 10,
	}

	// When
	err := company.Validate()

	// Then
	assert.Error(t, err, "Expected error for invalid Address")
}

func TestGivenCompany_WhenPWDEqualsEmployeeCount_ThenShouldReturnNoError(t *testing.T) {
	company := Company{
		CNPJ:                        "47960950000121",
		FantasyName:                 "Empresa Teste",
		CorporateName:               "Empresa Teste LTDA",
		Address:                     "Rua Teste, 123",
		EmployeeCount:               10,
		RequiredMinPWDEmployeeCount: 10, // Igual ao total
	}
	assert.NoError(t, company.Validate())
}

package domain

import (
	"company-service/pkg/utils"
	"errors"
	"time"
	"unicode/utf8"
)

type Company struct {
	ID                          string    `bson:"_id,omitempty" json:"id"`
	CNPJ                        string    `bson:"cnpj" json:"cnpj"`
	FantasyName                 string    `bson:"fantasy_name" json:"fantasy_name"`
	CorporateName               string    `bson:"corporate_name" json:"corporate_name"`
	Address                     string    `bson:"address" json:"address"`
	EmployeeCount               int       `bson:"employee_count" json:"employee_count"`
	RequiredMinPWDEmployeeCount int       `bson:"required_min_pwd_employee_count" json:"required_min_pwd_employee_count"`
	CreatedAt                   time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt                   time.Time `bson:"updated_at" json:"updated_at"`
}

func (c *Company) Validate() error {
	// Validação do CNPJ
	if err := c.validateCNPJ(); err != nil {
		return err
	}

	// Validação dos nomes
	if err := c.validateNames(); err != nil {
		return err
	}

	//Validação dos campos numéricos
	if err := c.validateNumbers(); err != nil {
		return err
	}

	// Validação dos campos obrigatórios
	if err := c.validateRequiredFields(); err != nil {
		return err
	}

	return nil
}

func (c *Company) validateCNPJ() error {
	if c.CNPJ == "" {
		return errors.New("CNPJ é obrigatório")
	}

	// Remove formatação para validação
	cleanCNPJ := utils.CleanCNPJ(c.CNPJ)
	if !utils.ValidCNPJ(cleanCNPJ) {
		return errors.New("CNPJ inválido")
	}

	c.CNPJ = cleanCNPJ

	return nil
}

func (c *Company) validateNames() error {
	// Validação do Nome Fantasia
	if c.FantasyName == "" {
		return errors.New("Nome Fantasia é obrigatório")
	}

	if utf8.RuneCountInString(c.FantasyName) < 2 {
		return errors.New("Nome Fantasia deve ter pelo menos 2 caracteres")
	}

	if utf8.RuneCountInString(c.FantasyName) > 100 {
		return errors.New("Nome Fantasia deve ter no máximo 100 caracteres")
	}

	// Validação da Razão Social
	if c.CorporateName == "" {
		return errors.New("Razão Social é obrigatória")
	}

	if utf8.RuneCountInString(c.CorporateName) < 5 {
		return errors.New("Razão Social deve ter pelo menos 5 caracteres")
	}

	if utf8.RuneCountInString(c.CorporateName) > 150 {
		return errors.New("Razão Social deve ter no máximo 150 caracteres")
	}

	return nil
}

func (c *Company) validateNumbers() error {
	// Validação da Quantidade de Funcionários
	if c.EmployeeCount < 0 {
		return errors.New("Quantidade de Funcionários não pode ser negativa")
	}

	// Validação da Quantidade Mínima de Funcionários PCD
	if c.RequiredMinPWDEmployeeCount < 0 {
		return errors.New("Quantidade Mínima de Funcionários PCD não pode ser negativa")
	}

	// Validação da Quantidade Mínima de Funcionários PCD
	if c.RequiredMinPWDEmployeeCount > c.EmployeeCount {
		return errors.New("Quantidade Mínima de Funcionários PCD não pode ser maior que a Quantidade de Funcionários")
	}

	// Validação básica baseada na Lei de Cotas para Pessoas com Deficiência (PCD) - Lei nº 8.213/1991 - Artigo 93.
	if c.EmployeeCount >= 100 && c.RequiredMinPWDEmployeeCount == 0 {
		return errors.New("Quantidade Mínima de Funcionários PCD é obrigatória para empresas com 100 ou mais funcionários")
	}

	return nil
}

func (c *Company) validateRequiredFields() error {
	// Validação dos campos obrigatórios
	if c.CNPJ == "" {
		return errors.New("CNPJ é obrigatório")
	}

	if c.FantasyName == "" {
		return errors.New("Nome Fantasia é obrigatório")
	}

	if c.CorporateName == "" {
		return errors.New("Razão Social é obrigatória")
	}

	if c.Address == "" {
		return errors.New("Endereço é obrigatório")
	}

	if c.EmployeeCount == 0 {
		return errors.New("Quantidade de Funcionários é obrigatória")
	}

	if c.RequiredMinPWDEmployeeCount == 0 {
		return errors.New("Quantidade Mínima de Funcionários PCD é obrigatória")
	}

	return nil
}

// BeforeUpdate hook para ser chamado antes de persistir
func (c *Company) BeforeCreate() {
	now := time.Now()
	if c.CreatedAt.IsZero() {
		c.CreatedAt = now
	}
	c.UpdatedAt = now
}

// BeforeUpdate hook para ser chamado antes de atualizar
func (c *Company) BeforeUpdate() {
	c.UpdatedAt = time.Now()
}

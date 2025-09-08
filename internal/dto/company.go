package dto

import (
	"company-service/internal/domain"
	"time"
)

// CreateCompanyRequest represents the request to create a new company.
type CreateCompanyRequest struct {
	CNPJ                        string `json:"cnpj" validate:"required,len=14"`
	FantasyName                 string `json:"fantasy_name" validate:"required"`
	CorporateName               string `json:"corporate_name" validate:"required"`
	Address                     string `json:"address" validate:"required"`
	EmployeeCount               int    `json:"employee_count" validate:"required"`
	RequiredMinPWDEmployeeCount int    `json:"required_min_pwd_employee_count" validate:"required"`
}

// UpdateCompanyRequest represents the request to update an existing company.
type UpdateCompanyRequest struct {
	CNPJ                        string `json:"cnpj" validate:"required,len=14"`
	FantasyName                 string `json:"fantasy_name" validate:"required"`
	CorporateName               string `json:"corporate_name" validate:"required"`
	Address                     string `json:"address" validate:"required"`
	EmployeeCount               int    `json:"employee_count" validate:"required"`
	RequiredMinPWDEmployeeCount int    `json:"required_min_pwd_employee_count" validate:"required"`
}

// CompanyResponse represents the response containing company details.
type CompanyResponse struct {
	ID                          string    `json:"id"`
	CNPJ                        string    `json:"cnpj"`
	FantasyName                 string    `json:"fantasy_name"`
	CorporateName               string    `json:"corporate_name"`
	Address                     string    `json:"address"`
	EmployeeCount               int       `json:"employee_count"`
	RequiredMinPWDEmployeeCount int       `json:"required_min_pwd_employee_count"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
}

// ToDomainCompany converts CreateCompanyRequest to domain.Company.
func ToDomainCompanyCreate(req *CreateCompanyRequest) *domain.Company {
	return &domain.Company{
		CNPJ:                        req.CNPJ,
		FantasyName:                 req.FantasyName,
		CorporateName:               req.CorporateName,
		Address:                     req.Address,
		EmployeeCount:               req.EmployeeCount,
		RequiredMinPWDEmployeeCount: req.RequiredMinPWDEmployeeCount,
	}
}

// ToDomainCompanyUpdate converts UpdateCompanyRequest to domain.Company.
func ToDomainCompanyUpdate(req *UpdateCompanyRequest, id string) *domain.Company {
	return &domain.Company{
		ID:                          id,
		CNPJ:                        req.CNPJ,
		FantasyName:                 req.FantasyName,
		CorporateName:               req.CorporateName,
		Address:                     req.Address,
		EmployeeCount:               req.EmployeeCount,
		RequiredMinPWDEmployeeCount: req.RequiredMinPWDEmployeeCount,
	}
}

func FromDomainCompany(company *domain.Company) *CompanyResponse {
	return &CompanyResponse{
		ID:                          company.ID,
		CNPJ:                        company.CNPJ,
		FantasyName:                 company.FantasyName,
		CorporateName:               company.CorporateName,
		Address:                     company.Address,
		EmployeeCount:               company.EmployeeCount,
		RequiredMinPWDEmployeeCount: company.RequiredMinPWDEmployeeCount,
		CreatedAt:                   company.CreatedAt,
		UpdatedAt:                   company.UpdatedAt,
	}
}

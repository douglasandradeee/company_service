package service

import (
	"company-service/internal/domain"
	"context"
)

// CompanyService define a interface para a camada de serviço
type CompanyService interface {
	CreateCompany(ctx context.Context, company *domain.Company) error
	GetCompany(ctx context.Context, id string) (*domain.Company, error)
	UpdateCompany(ctx context.Context, company *domain.Company) error
	DeleteCompany(ctx context.Context, id string) error
	ListCompanies(ctx context.Context, page, limit int) ([]*domain.Company, error)
}

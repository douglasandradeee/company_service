package repository

import (
	"company-service/internal/domain"
	"context"
)

// ompanyRepository define a interface para operações de persistência de empresas
type CompanyRepository interface {
	Create(ctx context.Context, company *domain.Company) error
	GetByID(ctx context.Context, id string) (*domain.Company, error)
	GetByCNPJ(ctx context.Context, cnpj string) (*domain.Company, error)
	Update(ctx context.Context, company *domain.Company) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*domain.Company, error)
	Count(ctx context.Context) (int64, error)
}

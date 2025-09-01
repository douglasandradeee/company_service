package domain

import "context"

// CompanyRepository é a interface que define os métodos para interagir/persistir com os dados da empresa
type CompanyRepository interface {
	Create(ctx context.Context) error
	GetByID(ctx context.Context, id string) (*Company, error)
	GetByCNPJ(ctx context.Context, cnpj string) (*Company, error)
	Update(ctx context.Context, company *Company) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int) ([]*Company, error)
	Count(ctx context.Context) (int64, error)
}

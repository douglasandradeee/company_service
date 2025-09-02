package messaging

import (
	"company-service/internal/domain"
	"context"
)

// MessageProducer define a produção de mensagens para eventos relacionados à empresa
type MessageProducer interface {
	SendCompanyCreated(ctx context.Context, company *domain.Company) error
	SendCompanyUpdated(ctx context.Context, company *domain.Company) error
	SendCompanyDeleted(ctx context.Context, company *domain.Company) error
	Close() error
}

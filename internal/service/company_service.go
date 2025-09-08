package service

import (
	"context"
	"fmt"
	"time"

	"company-service/internal/domain"
	"company-service/internal/messaging"
	"company-service/internal/repository"

	"go.uber.org/zap"
)

type companyService struct {
	repo            repository.CompanyRepository
	messageProducer messaging.MessageProducer
	logger          *zap.Logger
	maxRetries      int
	retryDelay      time.Duration
}

// NewCompanyService cria uma nova instância de CompanyService.
func NewCompanyService(repo repository.CompanyRepository, messageProducer messaging.MessageProducer, logger *zap.Logger) CompanyService {
	return &companyService{
		repo:            repo,
		messageProducer: messageProducer,
		logger:          logger,
		maxRetries:      3,               // Número máximo de tentativas
		retryDelay:      1 * time.Second, // Delay inicial entre tentativas
	}
}

// CreateCompany cria uma nova empresa.
func (s *companyService) CreateCompany(ctx context.Context, company *domain.Company) error {
	// Validação dos dados de entrada
	if err := company.Validate(); err != nil {
		return NewServiceError(err, "dados da empresa inválidos", "VALIDATION_ERROR")
	}

	// Verifica se o CNPJ já existe
	existing, err := s.repo.GetByCNPJ(ctx, company.CNPJ)
	if err != nil {
		return NewServiceError(err, "erro ao verificar CNPJ", "REPOSITORY_ERROR")
	}
	if existing != nil {
		return NewServiceError(ErrCNPJAlreadyExists, fmt.Sprintf("CNPJ %s já cadastrado", company.CNPJ), "CNPJ_CONFLICT")
	}

	// Hook para createdAT e updatedAT
	company.BeforeCreate()

	// Persiste empresa no repositório
	if err := s.repo.Create(ctx, company); err != nil {
		return NewServiceError(err, "erro ao criar empresa", "REPOSITORY_ERROR")
	}

	// Envia mensagem para o RabbitMQ com retry (async - não bloqueia)
	go s.sendCompanyCreatedMessage(context.Background(), company)

	s.logger.Info("Empresa criada com sucesso",
		zap.String("company_id", company.ID),
		zap.String("cnpj", company.CNPJ))

	return nil
}

// GetCompany busca uma empresa pelo ID.
func (s *companyService) GetCompany(ctx context.Context, id string) (*domain.Company, error) {
	if id == "" {
		return nil, NewServiceError(ErrInvalidCompanyData, "ID é obrigatório", "VALIDATION_ERROR")
	}

	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, NewServiceError(err, "erro ao buscar empresa", "REPOSITORY_ERROR")
	}
	if company == nil {
		return nil, NewServiceError(ErrCompanyNotFound, fmt.Sprintf("Empresa com ID %s não encontrada", id), "NOT_FOUND")
	}
	return company, nil
}

// UpdateCompany atualiza uma empresa existente.
func (s *companyService) UpdateCompany(ctx context.Context, company *domain.Company) (*domain.Company, error) {
	// Valida os dados de entrada
	if err := company.Validate(); err != nil {
		return nil, NewServiceError(err, "dados da empresa inválidos", "VALIDATION_ERROR")
	}

	// Verifica se a empresa existe
	existing, err := s.repo.GetByID(ctx, company.ID)
	if err != nil {
		return nil, NewServiceError(err, "erro ao buscar empresa", "REPOSITORY_ERROR")
	}
	if existing == nil {
		return nil, NewServiceError(ErrCompanyNotFound, fmt.Sprintf("Empresa com ID %s não encontrada", company.ID), "NOT_FOUND")
	}

	// Verifica se CNPJ foi alterado e se novo CNPJ já existe
	if existing.CNPJ != company.CNPJ {
		cnpjExists, err := s.repo.GetByCNPJ(ctx, company.CNPJ)
		if err != nil {
			return nil, NewServiceError(err, "erro ao verificar CNPJ", "REPOSITORY_ERROR")
		}
		if cnpjExists != nil {
			return nil, NewServiceError(ErrCNPJAlreadyExists, fmt.Sprintf("CNPJ %s já cadastrado", company.CNPJ), "CNPJ_CONFLICT")
		}
	}

	// Hook para updatedAT
	company.BeforeUpdate()

	// Persiste empresa no repositório
	updateCompany, err := s.repo.Update(ctx, company)
	if err != nil {
		return nil, NewServiceError(err, "erro ao atualizar empresa", "REPOSITORY_ERROR")
	}

	// Envia mensagem para o RabbitMQ com retry (async - não bloqueia)
	go s.sendCompanyUpdatedMessage(context.Background(), updateCompany)

	s.logger.Info("Empresa atualizada com sucesso",
		zap.String("company_id", updateCompany.ID))

	return updateCompany, nil
}

// DeleteCompany remove uma empresa.
func (s *companyService) DeleteCompany(ctx context.Context, id string) error {
	if id == "" {
		return NewServiceError(ErrInvalidCompanyData, "ID é obrigatório", "VALIDATION_ERROR")
	}

	company, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return NewServiceError(err, "erro ao buscar empresa", "REPOSITORY_ERROR")
	}
	if company == nil {
		return NewServiceError(ErrCompanyNotFound, fmt.Sprintf("Empresa com ID %s não encontrada", id), "NOT_FOUND")
	}

	// Remove empresa do repositório
	if err := s.repo.Delete(ctx, id); err != nil {
		return NewServiceError(err, "erro ao deletar empresa", "REPOSITORY_ERROR")
	}

	// Envia mensagem para o RabbitMQ com retry (async - não bloqueia)
	go s.sendCompanyDeletedMessage(context.Background(), company)

	s.logger.Info("Empresa removida com sucesso",
		zap.String("company_id", company.ID))

	return nil
}

// ListCompanies lista empresas com paginação.
func (s *companyService) ListCompanies(ctx context.Context, page int, limit int) ([]*domain.Company, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	companies, err := s.repo.List(ctx, page, limit)
	if err != nil {
		return nil, NewServiceError(err, "erro ao listar empresas", "REPOSITORY_ERROR")
	}

	return companies, nil
}

// sendWithRetry implementa mecanismo de retry para envio de mensagens
func (s *companyService) sendWithRetry(ctx context.Context, sendFunc func(ctx context.Context) error, operation string, companyID string) {
	for attempt := 1; attempt <= s.maxRetries; attempt++ {
		err := sendFunc(ctx)
		if err == nil {
			s.logger.Info("Mensagem enviada com sucesso",
				zap.String("operation", operation),
				zap.String("company_id", companyID))
			return
		}

		s.logger.Warn("Falha ao enviar mensagem, tentando novamente",
			zap.String("operation", operation),
			zap.String("company_id", companyID),
			zap.Int("attempt", attempt),
			zap.Error(err))

		if attempt < s.maxRetries {
			time.Sleep(time.Duration(attempt) * s.retryDelay)
		}
	}

	s.logger.Error("Falha ao enviar mensagem após todas as tentativas",
		zap.String("operation", operation),
		zap.String("company_id", companyID))
}

func (s *companyService) sendCompanyCreatedMessage(ctx context.Context, company *domain.Company) {
	s.sendWithRetry(ctx, func(ctx context.Context) error {
		return s.messageProducer.SendCompanyCreated(ctx, company)
	}, "company_created", company.ID)
}

func (s *companyService) sendCompanyUpdatedMessage(ctx context.Context, company *domain.Company) {
	s.sendWithRetry(ctx, func(ctx context.Context) error {
		return s.messageProducer.SendCompanyUpdated(ctx, company)
	}, "company_updated", company.ID)
}

func (s *companyService) sendCompanyDeletedMessage(ctx context.Context, company *domain.Company) {
	s.sendWithRetry(ctx, func(ctx context.Context) error {
		return s.messageProducer.SendCompanyDeleted(ctx, company)
	}, "company_deleted", company.ID)
}

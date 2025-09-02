package service

import "errors"

// Erros customizados do domínio
var (
	ErrCompanyNotFound    = errors.New("empresa não encontrada")
	ErrCNPJAlreadyExists  = errors.New("CNPJ já cadastrado")
	ErrInvalidCompanyData = errors.New("dados da empresa inválidos")
)

// ServiceError representa um erro na camada de serviço e encapsula erros com contexto adicional
type ServiceError struct {
	Err     error
	Message string
	Code    string
}

func (e *ServiceError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func NewServiceError(err error, message, code string) *ServiceError {
	return &ServiceError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}

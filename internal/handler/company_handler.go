package handler

import (
	"company-service/internal/domain"
	"company-service/internal/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type CompanyHandler struct {
	service service.CompanyService
	logger  *zap.Logger
}

func NewCompanyHandler(service service.CompanyService, logger *zap.Logger) *CompanyHandler {
	return &CompanyHandler{
		service: service,
		logger:  logger,
	}
}

// CreateCompanyHandler lida com a criação de uma nova empresa
func (h *CompanyHandler) CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to create company")

	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}

	if err := h.service.CreateCompany(r.Context(), &company); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(company); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}

// GetCompanyHandler lida com a busca de uma empresa por ID
func (h *CompanyHandler) GetCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.logger.Info("Received request to get company", zap.String("id", id))

	company, err := h.service.GetCompany(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(company); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}

// UpdateCompanyHandler lida com a atualização de uma empresa
func (h *CompanyHandler) UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.logger.Info("Received request to update company", zap.String("id", id))

	var company domain.Company
	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		h.logger.Error("Failed to decode request body", zap.Error(err))
		http.Error(w, `{"error": "Invalid JSON format"}`, http.StatusBadRequest)
		return
	}
	company.ID = id

	if err := h.service.UpdateCompany(r.Context(), &company); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(company); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}

// DeleteCompanyHandler lida com a exclusão de uma empresa
func (h *CompanyHandler) DeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	h.logger.Info("Received request to delete company", zap.String("id", id))

	if err := h.service.DeleteCompany(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListCompaniesHandler lida com a listagem de empresas com paginação
func (h *CompanyHandler) ListCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received request to list companies")

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	companies, err := h.service.ListCompanies(r.Context(), page, limit)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Create response with pagination info
	response := map[string]interface{}{
		"page":      page,
		"limit":     limit,
		"companies": companies,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode response", zap.Error(err))
	}
}

// handleServiceError trata os erros do service layer e retorna respostas HTTP apropriadas
func (h *CompanyHandler) handleServiceError(w http.ResponseWriter, err error) {
	if serviceErr, ok := err.(*service.ServiceError); ok {
		switch serviceErr.Code {
		case "VALIDATION_ERROR", "CNPJ_CONFLICT":
			h.logger.Warn("Validation error", zap.Error(err))
			http.Error(w, `{"error": "`+serviceErr.Error()+`"}`, http.StatusBadRequest)
		case "NOT_FOUND":
			h.logger.Warn("Resource not found", zap.Error(err))
			http.Error(w, `{"error": "`+serviceErr.Error()+`"}`, http.StatusNotFound)
		default:
			h.logger.Error("Service error", zap.Error(err))
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	h.logger.Error("Unexpected error", zap.Error(err))
	http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
}

// HealthCheckHandler fornece endpoint de health check
func (h *CompanyHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "service": "company-service"}`))
}

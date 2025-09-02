package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"company-service/internal/config"
	"company-service/internal/handler"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server struct {
	router *mux.Router
	logger *zap.Logger
	cfg    *config.Config
}

func NewServer(companyHandler *handler.CompanyHandler, logger *zap.Logger, cfg *config.Config) *Server {
	router := mux.NewRouter()

	// Configurar rotas
	router.HandleFunc("/companies", companyHandler.CreateCompanyHandler).Methods("POST")
	router.HandleFunc("/companies/{id}", companyHandler.GetCompanyHandler).Methods("GET")
	router.HandleFunc("/companies/{id}", companyHandler.UpdateCompanyHandler).Methods("PUT")
	router.HandleFunc("/companies/{id}", companyHandler.DeleteCompanyHandler).Methods("DELETE")
	router.HandleFunc("/companies", companyHandler.ListCompaniesHandler).Methods("GET")
	router.HandleFunc("/health", companyHandler.HealthCheckHandler).Methods("GET")

	// Middleware para logging
	router.Use(loggingMiddleware(logger))

	return &Server{
		router: router,
		logger: logger,
		cfg:    cfg,
	}
}

func (s *Server) Start() error {
	server := &http.Server{
		Addr:         ":" + s.cfg.ServerPort,
		Handler:      s.router,
		ReadTimeout:  parseDuration(s.cfg.ReadTimeout, 5*time.Second),
		WriteTimeout: parseDuration(s.cfg.WriteTimeout, 10*time.Second),
		IdleTimeout:  parseDuration(s.cfg.IdleTimeout, 60*time.Second),
	}

	// Canal para shutdown graceful
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		s.logger.Info("Shutting down server gracefully...")

		ctx, cancel := context.WithTimeout(context.Background(),
			parseDuration(s.cfg.ShutdownTimeout, 10*time.Second))
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			s.logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	s.logger.Info("Server starting", zap.String("address", server.Addr))
	return server.ListenAndServe()
}

func parseDuration(durationStr string, defaultDuration time.Duration) time.Duration {
	if duration, err := time.ParseDuration(durationStr); err == nil {
		return duration
	}
	return defaultDuration
}

// loggingMiddleware adiciona logging para todas as requests
func loggingMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Criar response writer customizado para capturar status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)

			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", wrapped.statusCode),
				zap.Duration("duration", duration),
				zap.String("ip", r.RemoteAddr),
			)
		})
	}
}

// responseWriter custom para capturar status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

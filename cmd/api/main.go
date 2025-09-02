package main

import (
	"company-service/internal/config"
	"company-service/pkg/logger"
	"log"

	"go.uber.org/zap"
)

func main() {

	// Carregando configurações
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Falha ao carregar as configurações: %v", err)
	}

	logger, err := logger.New(cfg.LogLevel)
	if err != nil {
		log.Fatal("Falha ao iniciar o logger")
	}
	defer logger.Sync()

	logger.Info("Configurações carregadas com sucesso")
	logger.Debug("Detalhes de Configuração:", zap.Any("config", cfg))
}

package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"

	"company-service/internal/config"
	"company-service/internal/handler"
	"company-service/internal/messaging/rabbitmq"
	"company-service/internal/repository/mongorepo"
	"company-service/internal/server"
	"company-service/internal/service"
)

func main() {
	// Carregar configuração
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Inicializar logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Failed to create logger:", err)
	}
	defer logger.Sync()

	logger.Info("Starting Company Service with configuration",
		zap.String("port", cfg.ServerPort),
		zap.String("mongo_db", cfg.MongoDB))

	// Conectar ao MongoDB
	mongoClient, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		logger.Fatal("Failed to connect to MongoDB", zap.Error(err))
	}
	defer mongoClient.Disconnect(context.Background())

	// Verificar conexão com MongoDB
	if err := mongoClient.Ping(context.Background(), nil); err != nil {
		logger.Fatal("Failed to ping MongoDB", zap.Error(err))
	}

	db := mongoClient.Database(cfg.MongoDB)

	// Inicializar repositório
	repo := mongorepo.NewCompanyRepositoryWithTimeout(
		db,
		cfg.MongoCollection,
		10*time.Second, // timeout de 10 segundos
	)

	logger.Info("MongoDB repository initialized",
		zap.String("database", cfg.MongoDB),
		zap.String("collection", cfg.MongoCollection))

	// Inicializar produtor RabbitMQ
	messageProducer, err := rabbitmq.NewProducer(cfg.RabbitMQURI, cfg.QueueName)
	if err != nil {
		logger.Fatal("Failed to create RabbitMQ producer", zap.Error(err))
	}
	defer messageProducer.Close()

	logger.Info("RabbitMQ producer initialized",
		zap.String("queue", cfg.QueueName))

	// Inicializar service
	companyService := service.NewCompanyService(repo, messageProducer, logger)

	// Inicializar handlers
	companyHandler := handler.NewCompanyHandler(companyService, logger)

	// Inicializar e iniciar servidor
	srv := server.NewServer(companyHandler, logger, cfg)

	logger.Info("Starting HTTP server", zap.String("address", ":"+cfg.ServerPort))
	if err := srv.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort      string `mapstructure:"SERVER_PORT"`
	MongoURI        string `mapstructure:"MONGO_URI"`
	MongoDB         string `mapstructure:"MONGO_DB"`
	MongoCollection string `mapstructure:"MONGO_COLLECTION"`
	RabbitMQURI     string `mapstructure:"RABBITMQ_URI"`
	QueueName       string `mapstructure:"QUEUE_NAME"`
	WebSocketPort   string `mapstructure:"WEBSOCKET_PORT"`
	LogLevel        string `mapstructure:"LOG_LEVEL"`
	ShutdownTimeout string `mapstructure:"SHUTDOWN_TIMEOUT"`
	ReadTimeout     string `mapstructure:"READ_TIMEOUT"`
	WriteTimeout    string `mapstructure:"WRITE_TIMEOUT"`
	IdleTimeout     string `mapstructure:"IDLE_TIMEOUT"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")

	// Configuração de valores padrão
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("MONGO_URI", "mongodb://localhost:27017")
	viper.SetDefault("MONGO_DB", "company_db")
	viper.SetDefault("MONGO_COLLECTION", "companies")
	viper.SetDefault("RABBITMQ_URI", "amqp://guest:guest@localhost:5672/")
	viper.SetDefault("QUEUE_NAME", "company_events")
	viper.SetDefault("WEBSOCKET_PORT", "8081")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("SHUTDOWN_TIMEOUT", "10s")
	viper.SetDefault("READ_TIMEOUT", "5s")
	viper.SetDefault("WRITE_TIMEOUT", "10s")
	viper.SetDefault("IDLE_TIMEOUT", "60s")

	// Lê variáveis de ambiente (tem precedência sobre o arquivo .env)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Configuração padrão será utilizada
		} else {
			return nil, fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("erro ao deserializar a configuração: %w", err)
	}
	return &config, nil
}

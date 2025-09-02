package rabbitmq

import (
	"company-service/internal/domain"
	"company-service/internal/messaging"
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

func NewProducer(uri, queueName string) (messaging.MessageProducer, error) {
	return &rabbitMQProducer{}, nil
}

func (p *rabbitMQProducer) SendCompanyCreated(ctx context.Context, company *domain.Company) error {
	return p.sendMessage(ctx, "company_created", "cadastro de EMPRESA"+company.FantasyName, company)
}

func (p *rabbitMQProducer) SendCompanyUpdated(ctx context.Context, company *domain.Company) error {
	return p.sendMessage(ctx, "company_updated", "Edição de EMPRESA"+company.FantasyName, company)
}

func (p *rabbitMQProducer) SendCompanyDeleted(ctx context.Context, company *domain.Company) error {
	return p.sendMessage(ctx, "company_deleted", "Exclusão de EMPRESA"+company.FantasyName, company)
}

func (p *rabbitMQProducer) sendMessage(ctx context.Context, eventType, message string, company *domain.Company) error {
	fmt.Println("Enviando mensagem para o RabbitMQ:", eventType, message, company.FantasyName)
	return nil
}

func (p *rabbitMQProducer) Close() error {
	return nil
}

package rabbitmq

import (
	"company-service/internal/domain"
	"company-service/internal/messaging"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitMQProducer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

func NewProducer(uri, queueName string) (messaging.MessageProducer, error) {
	// Create a new RabbitMQ connection
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// Create a new channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = channel.QueueDeclare(
		queueName,
		true,  // durable - sobrevive a reinicialização do broker
		false, // auto-delete quando não usado
		false, // exclusive - apenas esta conexão
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	log.Printf("Connected to RabbitMQ and declared queue: %s", queueName)

	return &rabbitMQProducer{
		conn:    conn,
		channel: channel,
		queue:   queueName,
	}, nil
}

func (p *rabbitMQProducer) SendCompanyCreated(ctx context.Context, company *domain.Company) error {
	return p.sendMessage(ctx, "company_created", "Cadastro de EMPRESA "+company.FantasyName, company)
}

func (p *rabbitMQProducer) SendCompanyUpdated(ctx context.Context, company *domain.Company) error {
	return p.sendMessage(ctx, "company_updated", "Edição da EMPRESA "+company.FantasyName, company)
}

func (p *rabbitMQProducer) SendCompanyDeleted(ctx context.Context, company *domain.Company) error {
	return p.sendMessage(ctx, "company_deleted", "Exclusão de EMPRESA"+company.FantasyName, company)
}

func (p *rabbitMQProducer) sendMessage(ctx context.Context, eventType, messageTxt string, company *domain.Company) error {

	message := map[string]interface{}{
		"operation":  messageTxt, // EX: "Cadastro de EMPRESA XXX"
		"company_id": company.ID,
		"cnpj":       company.CNPJ,
		"timestamp":  time.Now().Format(time.RFC3339),
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		"",      // exchange
		p.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Message sent: %s", messageTxt)
	return nil
}

func (p *rabbitMQProducer) Close() error {
	return nil
}

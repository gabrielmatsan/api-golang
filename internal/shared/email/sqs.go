package email

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/google/uuid"
)

type EmailMessage struct {
	To       string `json:"to" validate:"required,email"`
	Subject  string `json:"subject" validate:"required"`
	Body     string `json:"body" validate:"required"`
	Template string `json:"template" validate:"omitempty"`
}

type SqsProducer struct {
	client   *sqs.Client
	queueUrl string
}

func NewSQSClient() *sqs.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		panic("Erro ao carregar config AWS: " + err.Error())
	}

	return sqs.NewFromConfig(cfg)
}

func NewSQSProducer(client *sqs.Client, queueUrl string) (*SqsProducer, error) {
	if client == nil {
		return nil, fmt.Errorf("client SQS não pode ser nil")
	}

	if queueUrl == "" {
		return nil, fmt.Errorf("URL da fila não pode estar vazia")
	}

	return &SqsProducer{
		client:   client,
		queueUrl: queueUrl,
	}, nil
}

func (p *SqsProducer) SendEmailMessage(ctx context.Context, msg EmailMessage) error {
	messageBody, err := json.Marshal(msg)

	if err != nil {
		return fmt.Errorf("erro ao serializar mensagem: %w", err)
	}

	input := &sqs.SendMessageInput{
		QueueUrl:               aws.String(p.queueUrl),
		MessageBody:            aws.String(string(messageBody)),
		MessageGroupId:         aws.String("email-group"),
		MessageDeduplicationId: aws.String(uuid.New().String()),
	}

	result, err := p.client.SendMessage(ctx, input)

	if err != nil {
		return fmt.Errorf("erro ao enviar mensagem: %w", err)
	}

	log.Printf("Mensagem enviada com sucesso! MessageId: %s", *result.MessageId)
	return nil
}

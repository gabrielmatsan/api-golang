package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// EmailMessage representa a estrutura da mensagem de e-mail (igual ao shared)
type EmailMessage struct {
	To       string `json:"to" validate:"required,email"`
	Subject  string `json:"subject" validate:"required"`
	Body     string `json:"body" validate:"required"`
	Template string `json:"template" validate:"omitempty"`
}

// EmailConsumer representa o consumer SQS para e-mails
type EmailConsumer struct {
	sesClient *ses.Client
	fromEmail string
}

// NewEmailConsumer cria uma nova instância do consumer de e-mail
func NewEmailConsumer() (*EmailConsumer, error) {
	// Carregar configuração AWS
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, err
	}

	// Obter e-mail padrão do remetente
	fromEmail := os.Getenv("DEFAULT_FROM_EMAIL")
	if fromEmail == "" {
		log.Printf("AVISO: DEFAULT_FROM_EMAIL não configurado")
	}

	return &EmailConsumer{
		sesClient: ses.NewFromConfig(cfg),
		fromEmail: fromEmail,
	}, nil
}

// ProcessMessage processa uma mensagem individual do SQS
func (c *EmailConsumer) ProcessMessage(ctx context.Context, messageBody string) error {
	// Parse da mensagem JSON
	var emailMsg EmailMessage
	if err := json.Unmarshal([]byte(messageBody), &emailMsg); err != nil {
		log.Printf("Erro ao fazer parse da mensagem: %v", err)
		return err
	}

	// Validar campos obrigatórios
	if emailMsg.To == "" || emailMsg.Subject == "" || emailMsg.Body == "" {
		log.Printf("Campos obrigatórios faltando na mensagem")
		return nil // Não retorna erro para não reprocessar mensagem inválida
	}

	// Enviar e-mail via SES
	if err := c.sendEmail(ctx, emailMsg); err != nil {
		log.Printf("Erro ao enviar e-mail para %s: %v", emailMsg.To, err)
		return err // Retorna erro para que a mensagem volte para a fila
	}

	log.Printf("E-mail enviado com sucesso para %s", emailMsg.To)
	return nil
}

// sendEmail envia um e-mail via Amazon SES
func (c *EmailConsumer) sendEmail(ctx context.Context, emailMsg EmailMessage) error {
	// Usar e-mail padrão se não especificado
	fromEmail := c.fromEmail
	if fromEmail == "" {
		fromEmail = "noreply@example.com" // Fallback
	}

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{emailMsg.To},
		},
		Message: &types.Message{
			Body: &types.Body{
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(emailMsg.Body),
				},
			},
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(emailMsg.Subject),
			},
		},
		Source: aws.String(fromEmail),
	}

	_, err := c.sesClient.SendEmail(ctx, input)
	return err
}

// handler é o ponto de entrada da função Lambda
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Criar consumer de e-mail
	consumer, err := NewEmailConsumer()
	if err != nil {
		log.Printf("Erro ao criar consumer de e-mail: %v", err)
		return err
	}

	// Processar cada mensagem do lote SQS
	for _, record := range sqsEvent.Records {
		log.Printf("Processando mensagem %s da fila %s", record.MessageId, record.EventSourceARN)

		if err := consumer.ProcessMessage(ctx, record.Body); err != nil {
			log.Printf("Erro ao processar mensagem %s: %v", record.MessageId, err)
			return err // Falha em qualquer mensagem causa reprocessamento do lote
		}
	}

	return nil
}

func main() {
	lambda.Start(handler)
}

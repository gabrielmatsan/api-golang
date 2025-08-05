package templates

import (
	"context"
	"fmt"
	"log"

	"github.com/gabrielmatsan/teste-api/internal/shared/email"
	"github.com/gabrielmatsan/teste-api/internal/shared/singlaton"
	"github.com/gabrielmatsan/teste-api/internal/user/model"
)

func SendWelcomeEmail(user model.User) {
	// Pega o SQS producer do singleton (ou injete via DI)
	producer := singlaton.GetSQSProducer() // Você precisa implementar isso

	// Cria mensagem de email
	emailMsg := email.EmailMessage{
		To:       user.Email,
		Subject:  "Bem-vindo à nossa plataforma!",
		Body:     fmt.Sprintf("Olá %s %s, bem-vindo! Sua conta foi criada com sucesso.", user.FirstName, user.LastName),
		Template: "welcome",
	}

	// Envia para a fila SQS
	if err := producer.SendEmailMessage(context.Background(), emailMsg); err != nil {
		log.Printf("❌ Erro ao enviar email de boas-vindas para %s: %v", user.Email, err)
	} else {
		log.Printf("✅ Email de boas-vindas enviado para a fila: %s", user.Email)
	}
}

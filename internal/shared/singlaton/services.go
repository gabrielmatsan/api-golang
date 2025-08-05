package singlaton

import (
	"os"
	"sync"

	"github.com/gabrielmatsan/teste-api/internal/shared/email"
)

var (
	sqsProducer *email.SqsProducer
	sqsOnce     sync.Once
)

func GetSQSProducer() *email.SqsProducer {
	sqsOnce.Do(func() {
		queueURL := os.Getenv("SQS_URL")
		if queueURL == "" {
			// Fallback para URL hardcoded ou panic
			queueURL = "https://sqs.us-east-1.amazonaws.com/123456789/pub-sub-golang.fifo"
		}

		// Inicializa o cliente SQS (supondo que você tenha uma função para isso)
		sqsClient := email.NewSQSClient() // Implemente esta função conforme necessário

		producer, err := email.NewSQSProducer(sqsClient, queueURL)
		if err != nil {
			panic("Erro ao inicializar SQS Producer: " + err.Error())
		}
		sqsProducer = producer
	})
	return sqsProducer
}

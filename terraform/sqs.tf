module "sqs" {
  source  = "terraform-aws-modules/sqs/aws"
  version = "5.0.0"


  # Fifo Configuration
  name                        = "pub-sub-golang"
  fifo_queue                  = true
  content_based_deduplication = true


  create_dlq = true
  dlq_name   = "pub-sub-golang-dlq"

  # Time Configuration
  visibility_timeout_seconds = 300     # 5 minutos para processar
  message_retention_seconds  = 1209600 # 14 dias (máximo)
  receive_wait_time_seconds  = 20      # Long polling
  delay_seconds              = 0       # Sem delay por padrão

  max_message_size = 262144 # 256KB (máximo)

  create_queue_policy = true
  queue_policy_statements = {
    send_receive = {
      sid = "AllowSendReceive"
      actions = [
        "sqs:SendMessage",
        "sqs:ReceiveMessage",
        "sqs:DeleteMessage",
        "sqs:GetQueueAttributes"
      ]
      principals = [
        {
          type        = "AWS"
          identifiers = ["arn:aws:iam::${data.aws_caller_identity.current.account_id}:root"]
        }
      ]
    }
  }

  redrive_policy = {
    maxReceiveCount = 5
  }

  tags = {
    IAC = "true"
  }
}


output "sqs_queue_id" {
  description = "URL da fila SQS"
  value       = module.sqs.queue_id
}

output "sqs_queue_arn" {
  description = "ARN da fila SQS"
  value       = module.sqs.queue_arn
}

output "sqs_queue_name" {
  description = "Nome da fila SQS"
  value       = module.sqs.queue_name
}

output "sqs_dlq_id" {
  description = "URL da DLQ"
  value       = module.sqs.dead_letter_queue_id
}

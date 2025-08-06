resource "aws_ses_email_identity" "sender" {
  email = "gabrielmatsan@hotmail.com"
}

# Configuração de regras de envio
resource "aws_ses_configuration_set" "email_config" {
  name = "email-consumer-config"

  delivery_options {
    tls_policy = "Require"
  }

  reputation_metrics_enabled = true
}

# Output para facilitar o uso
output "ses_sender_email" {
  description = "E-mail verificado no SES"
  value       = aws_ses_email_identity.sender.email
} 
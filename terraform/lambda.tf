resource "aws_iam_role" "lambda_role" {
  name = "email-processor-lambda-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "lambda.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_basic" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
  role       = aws_iam_role.lambda_role.name
}

data "archive_file" "lambda_zip" {
  type        = "zip"
  source_file = "${path.module}/../bootstrap"
  output_path = "${path.module}/../lambda.zip"
}

resource "aws_lambda_function" "email_consumer" {
  function_name = "email-consumer"
  role          = aws_iam_role.lambda_role.arn

  filename         = data.archive_file.lambda_zip.output_path
  source_code_hash = data.archive_file.lambda_zip.output_base64sha256

  handler = "bootstrap"
  runtime = "provided.al2"

  environment {
    variables = {
      ENV                = "dev"
      DEFAULT_FROM_EMAIL = aws_ses_email_identity.sender.email
    }
  }

  tags = {
    IAC = "true"
  }
}

# Política para SQS
resource "aws_iam_policy" "sqs_policy" {
  name        = "email-consumer-sqs-policy"
  description = "Policy for SQS access from Lambda"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "sqs:ReceiveMessage",
          "sqs:DeleteMessage",
          "sqs:GetQueueAttributes"
        ],
        Effect   = "Allow",
        Resource = module.sqs.queue_arn
      }
    ]
  })
}

# Política para SES (envio de e-mails)
resource "aws_iam_policy" "ses_policy" {
  name        = "email-consumer-ses-policy"
  description = "Policy for SES access from Lambda"

  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Action = [
          "ses:SendEmail",
          "ses:SendRawEmail"
        ],
        Effect   = "Allow",
        Resource = "*"
      }
    ]
  })
}
# Anexar política SQS à role
resource "aws_iam_role_policy_attachment" "sqs_attachment" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.sqs_policy.arn
}

# Anexar política SES à role
resource "aws_iam_role_policy_attachment" "ses_attachment" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = aws_iam_policy.ses_policy.arn
}

# Mapeamento de eventos SQS -> Lambda
resource "aws_lambda_event_source_mapping" "email_consumer_mapping" {
  event_source_arn = module.sqs.queue_arn
  function_name    = aws_lambda_function.email_consumer.arn
  batch_size       = 1
}
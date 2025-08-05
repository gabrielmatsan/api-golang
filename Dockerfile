FROM golang:1.24 AS base


# Baixando as dependências
FROM base AS dependencies
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download

# Buildando a aplicação
FROM dependencies AS builder
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd


FROM alpine:latest AS runner
# Criando o grupo e o usuário nao root
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /usr/src/app/main .
RUN chown appuser:appgroup main && \
    chmod +x main

USER appuser

EXPOSE 8080

CMD ["./main"]

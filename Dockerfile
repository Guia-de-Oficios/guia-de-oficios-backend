FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /guia-oficios-api ./cmd/api

FROM alpine:3.20
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /guia-oficios-api /app/guia-oficios-api
COPY migrations /app/migrations
RUN mkdir -p /app/uploads

EXPOSE 8080
ENV PORT=8080
ENV UPLOAD_DIR=/app/uploads

CMD ["/app/guia-oficios-api"]

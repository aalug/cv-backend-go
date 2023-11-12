FROM golang:1.21-alpine3.18 AS builder

WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Set gin to production
ENV GIN_MODE=release

RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.1/migrate.linux-amd64.tar.gz  | tar xvz

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .

COPY app.env .
COPY wait-for.sh .
COPY start.sh .

COPY internal/db/migrations ./migrations

EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT ["/app/start.sh"]
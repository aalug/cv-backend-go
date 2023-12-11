FROM golang:1.21-alpine3.18 AS builder

WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

# Install necessary tools (curl and wget)
RUN apk add --no-cache curl wget

# Download wait-for
RUN wget -qO wait-for.sh https://raw.githubusercontent.com/eficode/wait-for/v2.2.3/wait-for \
    && chmod +x wait-for.sh

RUN wget -qO start.sh https://gist.githubusercontent.com/aalug/deae8a5de108fc84dc8e52c74bc92cd3/raw/b8886c21f9d337cd3c8f97a46c2a83afc09a8262/start.sh \
    && chmod +x start.sh

# Download golang-migrate and place it in /app/migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.1/migrate.linux-amd64.tar.gz  | tar xvz

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/wait-for.sh .
COPY --from=builder /app/start.sh .
COPY --from=builder /app/migrate .

COPY app.env .

# Set gin to production
ENV GIN_MODE=release

# Set to production, if used locally - set to false
ENV PRODUCTION=true

COPY internal/db/migrations ./migrations

EXPOSE 8080
CMD ["./migrate", "-path", "migrations", "-database", "postgres://user:password@host:port/database", "up"]
ENTRYPOINT ["/app/start.sh"]
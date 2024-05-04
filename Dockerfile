# build stage
FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

# run stage
From alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY  config/config.json /app/config/
COPY db/migration ./migration
COPY start.sh .
RUN chmod +x start.sh

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
# build stage
FROM golang:1.22.2-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# run stage
From alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .

RUN mkdir config
COPY  config/config.json /app/config/

EXPOSE 8080
CMD [ "/app/main" ]
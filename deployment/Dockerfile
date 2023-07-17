FROM golang:1.20 AS builder
WORKDIR /app
COPY src/ .
ENV CGO_ENABLED=0
RUN go build -o main .

FROM alpine:latest
WORKDIR /app
RUN mkdir ssl
COPY ./ssl/cert.pem ssl/
COPY ./ssl/key.pem ssl/
COPY --from=builder /app/main .

ENV ALLOW_ORIGINS=https://localhost:4200
CMD ["./main"]
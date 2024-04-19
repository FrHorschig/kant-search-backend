FROM golang:1.21 AS builder
WORKDIR /app
COPY src/kant-search-backend .
RUN mkdir ssl
COPY ./ssl/cert.pem ssl/
COPY ./ssl/key.pem ssl/

ENV ALLOW_ORIGINS=https://localhost:4200
CMD ["./kant-search-backend"]
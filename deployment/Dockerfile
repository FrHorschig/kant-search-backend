FROM golang:1.24

LABEL org.opencontainers.image.title="kant-search-backend"
LABEL org.opencontainers.image.source="https://github.com/FrHorschig/kant-search-backend"
LABEL org.opencontainers.image.licenses="AGPL-3.0 license"

WORKDIR /app
COPY src/kant-search-backend .
RUN chmod +x kant-search-backend

CMD ["./kant-search-backend"]
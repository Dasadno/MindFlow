FROM node:22-alpine AS frontend-builder
WORKDIR /web
COPY client/package*.json ./
RUN npm install
COPY client/ .
RUN npm run build

# --- Сборка бэкенда (Go 1.26) ---
FROM golang:1.26-bookworm AS backend-builder

# Ставим зависимости для сборки CGO
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libc6-dev \
    libsqlite3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /web/dist ./web/dist

RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -o main ./server/cmd/server/main.go

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libsqlite3-0 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=backend-builder /app/main .
COPY --from=frontend-builder /web/dist ./web/dist

RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./main"]

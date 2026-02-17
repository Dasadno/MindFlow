FROM node:22-alpine AS frontend-builder
WORKDIR /web
COPY web/package*.json ./
RUN npm install
COPY web/ .
RUN npm run build

FROM golang:1.24-alpine AS backend-builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=frontend-builder /web/dist ./web/dist

RUN CGO_ENABLED=1 GOOS=linux go build -trimpath -o main .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates sqlite-libs

WORKDIR /root/
COPY --from=backend-builder /app/main .

RUN mkdir /data

EXPOSE 8080
CMD ["./main"]

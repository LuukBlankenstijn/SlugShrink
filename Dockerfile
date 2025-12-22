FROM node:24-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN corepack enable
RUN pnpm install
COPY frontend/ ./
RUN pnpm build

FROM golang:1.25-alpine AS backend-builder
WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./
RUN go mod download

COPY ./backend .
COPY --from=frontend-builder /app/frontend/build ./frontend/build

RUN CGO_ENABLED=0 go build -tags prod -ldflags="-w -s" -o main ./cmd/server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=backend-builder /app/main .

ENV DB_PATH=/data/sqlite.db
VOLUME /data

EXPOSE 8080
EXPOSE 8081
CMD ["./main"]

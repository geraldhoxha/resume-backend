FROM golang:1.24.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
    && CGO_ENABLED=0 GOOS=linux go build -o ./main



FROM alpine:latest

WORKDIR /app

COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY --from=builder /app/sql /app/sql
COPY --from=builder /app/main /app/wait-for-it.sh /app/run-all.sh /app/

RUN apk update && apk upgrade && apk add bash && apk add dos2unix \
    && chmod +x /app/main /app/wait-for-it.sh /app/run-all.sh \
    && find ./ type f -name "*.sh" -print0 | xargs -0 dos2unix \
    && apk del dos2unix

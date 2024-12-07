FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN CGO_ENABLED=0 GOOS=linux go build -o ./main
RUN chmod +x ./wait-for-it.sh ./run-all.sh

EXPOSE 8080

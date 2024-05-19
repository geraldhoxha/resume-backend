run:
	go run main.go

run-server:
	go run server.go

build:
	go build .

clean:
	rm -rf ./build

gen:
	go mod tidy
	go generate ./..
	go run github.com/99designs/gqlgen generate	

dbconnect:
	docker exec -it pg_db psql -U pguser -W pg_db

dbinit:
	migrate -path sql -database "postgres://pguser:SECRET@127.0.0.1:5432/pg_db?sslmode=disable" up


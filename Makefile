TOOLS_FILE:=./tools/tools.go
GOLANGCI:=2.1.5
run:
	go run main.go

run-server:
	go run server.go

lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI)
	golangci-lint run

build:
	go build .

clean:
	rm -rf ./build

gen: create_tools
	go mod tidy
	go generate ./...
	go run github.com/99designs/gqlgen generate
	rm $(TOOLS_FILE)

dbconnect:
	docker exec -it pg_db psql -U pguser -W pg_db

dbinit:
	migrate -path sql -database "postgres://pguser:SECRET@pghost:5432/pg_db?sslmode=disable" up

dockemigrate:
	docker run -v {{ migration dir }}:/migrations --network host migrate/migrate -path=/migrations/ -database postgres://pghost:5432/database up

create_tools:
	@echo 'package tools' > $(TOOLS_FILE)
	@echo 'import (' >> $(TOOLS_FILE)
	@echo ' _ "github.com/99designs/v1.62.2gqlgen"' >> $(TOOLS_FILE)
	@echo ')' >> $(TOOLS_FILE)

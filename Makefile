TOOLS_FILE:=./tools/tools.go

run:
	go run main.go

run-server:
	go run server.go

build:
	go build .

clean:
	rm -rf ./build

gen: create_tools
	go mod tidy
	go generate ./..
	go run github.com/99designs/gqlgen generate
	rm $(TOOLS_FILE)

dbconnect:
	docker exec -it pg_db psql -U pguser -W pg_db

dbinit:
	migrate -path sql -database "postgres://pguser:SECRET@127.0.0.1:5432/pg_db?sslmode=disable" up

create_tools:
	@echo 'package tools' > $(TOOLS_FILE)
	@echo 'import (' >> $(TOOLS_FILE)
	@echo ' _ "github.com/99designs/gqlgen"' >> $(TOOLS_FILE)
	@echo ')' >> $(TOOLS_FILE)

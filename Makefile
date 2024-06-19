build:
	cd api/cmd && go build main.go
lint:
	cd api && golangci-lint run
run: build
	cd api/cmd && ./main ../.env
goose-up:
	cd api/migrations && goose postgres "user=postgres password=postgres host=localhost port=5432 dbname=keep_it sslmode=disable" up
goose-down:
	cd api/migrations && goose postgres "user=postgres password=postgres host=localhost port=5432 dbname=keep_it sslmode=disable" down
goose-down-all:
	cd api/migrations && goose postgres "user=postgres password=postgres host=localhost port=5432 dbname=keep_it sslmode=disable" down-to 0
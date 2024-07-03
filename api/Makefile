build:
	cd cmd && go build main.go
lint:
	golangci-lint run
run: build
	cd cmd && ./main ../.env
goose-up:
	cd migrations && goose postgres "user=postgres password=postgres host=localhost port=5432 dbname=keep_it sslmode=disable" up
goose-down:
	cd migrations && goose postgres "user=postgres password=postgres host=localhost port=5432 dbname=keep_it sslmode=disable" down
goose-down-all:
	cd migrations && goose postgres "user=postgres password=postgres host=localhost port=5432 dbname=keep_it sslmode=disable" down-to 0
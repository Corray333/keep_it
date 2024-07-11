build:
	cd cmd && go build main.go
lint:
	golangci-lint run
run: build
	cd cmd && ./main ../.env

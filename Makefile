generate_proto:
	buf generate

tidy:
	go mod tidy

lint: tidy
	go vet ./...

format: tidy
	go fmt ./...

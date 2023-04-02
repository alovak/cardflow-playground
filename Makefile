build:
	mkdir -p bin
	go build -o bin/issuer -v ./cmd/issuer
	go build -o bin/acquirer -v ./cmd/acquirer


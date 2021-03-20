run:
	go run cmd/backend/main.go

build:
	go build -o bin/backend cmd/backend/main.go
	go build -o bin/cluster-client cmd/cluster-client/main.go
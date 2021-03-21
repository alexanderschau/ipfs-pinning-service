run-backend:
	go run cmd/backend/main.go

run-cluster-client:
	go run cmd/cluster-client/main.go

build:
	make build-backend
	make build-cluster-client

build-backend:
	go build -o bin/backend cmd/backend/main.go

build-cluster-client:
	go build -o bin/cluster-client cmd/cluster-client/main.go

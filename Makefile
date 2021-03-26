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

build-images:
	docker build . -f Dockerfile.backend -t $IMAGE_BACKEND_NAME:$IMAGE_VERSION
	docker build . -f Dockerfile.clusterClient -t $IMAGE_CLUSTER_CLIENT_NAME:$IMAGE_VERSION

push-images:
	docker push $IMAGE_BACKEND_NAME:$IMAGE_VERSION
	docker push $IMAGE_CLUSTER_CLIENT_NAME:$IMAGE_VERSION
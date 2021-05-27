run:
	go run cmd/local/main.go
build-oapi:
	openapi-generator-cli generate -g go -i https://raw.githubusercontent.com/ipfs/pinning-services-api-spec/master/ipfs-pinning-service.yaml -o openapi
	rm openapi/go.mod openapi/go.sum
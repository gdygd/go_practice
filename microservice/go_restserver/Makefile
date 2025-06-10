.PHONY: test server build
.EXPORT_ALL_VARIABLES:

test:
	go test -v -cover -short ./...

server:
	cd ./cli && make build
	cd ./mpserver && make build
	cd ./restserver && make build
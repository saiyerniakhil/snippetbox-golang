BINARY_NAME=snippetbox
GOPATH=/Users/madakhil/go

clean:
	rm -f ./bin/${BINARY_NAME}*


start-delve:
	$(GOPATH)/bin/dlv exec ./bin/$(BINARY_NAME)-debug --listen=127.0.0.1:2345 --headless=true --api-version=2 --accept-multiclient --continue --log -- 

build-debug: clean
	CGO_ENABLED=0 go build -gcflags=all="-N -l" -o bin/$(BINARY_NAME)-debug ./cmd/web/

codegen:
	sqlc generate .

start:
	$(GOPATH)/bin/air
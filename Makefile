.PHONY: mockgen
mockgen:
	mockgen \
		-source=internal/quote/repository.go \
		-destination=internal/quote/mock/repository.go \
		-package=mock
	mockgen \
		-source=internal/transport/interfaces.go \
		-destination=internal/transport/mock/interfaces.go \
		-package=mock

.PHONY: test
test:
	go test -v -count=1 ./... -coverprofile=cover.out

.PHONY: server
server:
	go run ./cmd/server

.PHONY: client
client:
	go run ./cmd/client

.PHONY: load
load:
	go run ./test/load

.PHONY: build-server
build-server:
	mkdir -p ./bin && \
	go build -o ./bin/server ./cmd/server

.PHONY: build-client
build-client:
	mkdir -p ./bin && \
	go build -o ./bin/client ./cmd/client

.PHONY: docker-build-server
docker-build-server:
	docker build -t powow-server:latest -f Dockerfile.server .

.PHONY: docker-build-client
docker-build-client:
	docker build -t powow-client:latest -f Dockerfile.client .

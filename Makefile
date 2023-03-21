# go build
GO_BUILD=go build -ldflags -s -v
BINARY_NAME=das_account_indexer_server

# update
update:
	go mod tidy

# linux
indexer_linux:
	export GOOS=linux
	export GOARCH=amd64
	$(GO_BUILD) -o $(BINARY_NAME) cmd/main.go
	mkdir -p bin/linux
	mv $(BINARY_NAME) bin/linux/
	@echo "build $(BINARY_NAME) successfully."

# mac
indexer_mac:
	export GOOS=darwin
	export GOARCH=amd64
	$(GO_BUILD) -o $(BINARY_NAME) cmd/main.go
	mkdir -p bin/mac
	mv $(BINARY_NAME) bin/mac/
	@echo "build $(BINARY_NAME) successfully."

# win
indexer_win: BINARY_NAME=das_account_indexer_server.exe
indexer_win:
	export GOOS=windows
	export GOARCH=amd64
	$(GO_BUILD) -o $(BINARY_NAME) cmd/main.go
	mkdir -p bin/win
	mv $(BINARY_NAME) bin/win/
	@echo "build $(BINARY_NAME) successfully."

docker:
	docker build --network host -t dotbitteam/das-account-indexer:latest .

docker-publish:
	docker image push dotbitteam/das-account-indexer:latest
# default
default: indexer_linux

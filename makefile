APP=camp
BIN=./bin/$(APP)
MAIN=cmd/camp/main.go

run: build
	@$(BIN)

build:
	@go build -o $(BIN) $(MAIN)

dev:
	@air

tidy:
	@go mod tidy

fmt:
	@go fmt ./...

clean:
	@rm -rf bin/* tmp/*


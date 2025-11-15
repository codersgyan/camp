run: build
	@./bin/camp

#  install---> go install github.com/air-verse/air@latest
dev:
	@command -v air >/dev/null 2>&1
	@air

build:
	@go build -o ./bin/camp cmd/camp/main.go
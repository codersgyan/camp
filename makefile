run: build
	@./bin/camp

#  install---> go install github.com/air-verse/air@latest
dev:
	@PATH=$$HOME/go/bin:$$PATH air

build:
	@go build -o ./bin/camp cmd/camp/main.go
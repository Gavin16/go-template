# variable definition
GO_CMD = go

.PHONY: build
build: clean swag-init
all: build run

clean:
	$(GO_CMD) clean
	rm -rf docs/*

swag-init:
	$(GO_CMD) run github.com/swaggo/swag/cmd/swag init

run:
	$(GO_CMD) run main.go


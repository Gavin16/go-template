# variable definition
GO_CMD = go
SWAG_CMD = github.com/swaggo/swag/cmd/swag


.PHONY: all clean swag-install swag-init tidy compile
init: clean swag-install swag-init tidy
build: clean swag-install swag-init tidy compile
all: clean swag-install swag-init tidy compile run


clean:
	@echo "Cleaning up ..."
	rm -rf docs/*

swag-install:
	$(GO_CMD) get -u $(SWAG_CMD)

swag-init:
	@echo "Generating API documentation..."
	$(GO_CMD) run github.com/swaggo/swag/cmd/swag init

tidy:
	$(GO_CMD) mod tidy

compile:
	$(GO_CMD) build -o ./go-template

run:
	./go-template
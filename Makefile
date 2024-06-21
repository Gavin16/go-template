# variable definition
GO_CMD = go
SWAG_CMD = github.com/swaggo/swag/cmd/swag


.PHONY: all clean tidy swag-install swag-init run
all: clean tidy swag-install swag-init run
build: clean tidy swag-install swag-init

clean:
	@echo "Cleaning up ..."
	rm -rf docs/*

tidy:
	$(GO_CMD) mod tidy

swag-install:
	$(GO_CMD) get -u $(SWAG_CMD)

swag-init:
	@echo "Generating API documentation..."
	$(GO_CMD) run github.com/swaggo/swag/cmd/swag init

#compile:
#	$(GO_CMD) build -o ./

run:
	$(GO_CMD) run main.go
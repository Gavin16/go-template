# variable definition
GO_CMD = go
SWAG_CMD = github.com/swaggo/swag/cmd/swag


.PHONY: all clean swag-install swag-init
all: clean swag-install swag-init
#build: clean tidy swag-install swag-init

clean:
	@echo "Cleaning up ..."
	rm -rf docs/*

swag-install:
	$(GO_CMD) get -u $(SWAG_CMD)

swag-init:
	@echo "Generating API documentation..."
	$(GO_CMD) run github.com/swaggo/swag/cmd/swag init



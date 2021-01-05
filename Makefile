APP         ?= xpath-linter
GOOS        ?= linux
VERSION     ?= 0.1.0
BUILD       ?= 0
BUILD_TIME  ?= $(shell date -u '+%F_%T%z')
COMMIT      ?= $(shell git rev-parse --short HEAD)
RELEASE_DIR ?= release
MODULE_NAME ?= githib.com/bop0hz/xpath-linter
LDFLAGS     ?= '-X $(MODULE_NAME)/version.builtAt=$(BUILD_TIME) \
				-X $(MODULE_NAME)/version.build=$(BUILD) \
				-X $(MODULE_NAME)/version.version=$(VERSION)'
CONFIGS     ?= configs

all: clean build test run

.PHONY: build
build: clean
	@echo "Building the $(APP)..\n"
	go build -v -o $(RELEASE_DIR)/$(APP) \
		-ldflags $(LDFLAGS) \
		main.go
	@cp -r ./examples $(RELEASE_DIR)/
	@cp ./config.yaml $(RELEASE_DIR)/

.PHONY: clean
clean:
	@echo "Cleanup $(RELEASE_DIR) dir.. \n"
	@rm -rf $(RELEASE_DIR)

.PHONY: run
run:
	@echo "Examples:"
	@echo "[CI mode]"
	find ./examples -type f -exec $(RELEASE_DIR)/$(APP) -ci {} \;
	@echo "[Adhoc mode]"
	./$(RELEASE_DIR)/$(APP) -must=false -contain '//node[text()="1"]' ./examples/multinodes.xml
	./$(RELEASE_DIR)/$(APP) -contain '//node[text()="4"]' ./examples/multinodes.xml

.PHONY: test
test:
	go test -v ./lint

.PHONY: generate run test

# Variables
OPENAPI_SPEC := api-spec/openapi.yaml
GENERATED_DIR := generated-server
TARGET_DIR := pkg/infra/httpinfragenerated

CUR_DIR := $(CURDIR)
CONFIG_PATH := $(CUR_DIR)

TEST_PATH := $(CUR_DIR)/test

# Generate server stubs and move to infrastructure layer
generate:
	@echo "Generating server code from OpenAPI spec..."
	openapi-generator generate \
		-i $(OPENAPI_SPEC) \
		-g go-server \
		-o $(GENERATED_DIR) \
		--package-name httpinfragenerated \
		--additional-properties=apiPackage=httpinfragenerated,modelPackage=httpinfragenerated

	@echo "Moving generated code to infrastructure layer..."
	rm -rf $(TARGET_DIR)
	mkdir -p $(TARGET_DIR)
	cp -r $(GENERATED_DIR)/go/* $(TARGET_DIR)/
	rm -rf $(GENERATED_DIR)

	@echo "Formatting generated code..."
	go fmt ./$(TARGET_DIR)/...

# Run the application
run:
	@echo "Starting server..."
	go run main.go start --configPath $(CONFIG_PATH)

test:
	go test -v $(TEST_PATH)/...


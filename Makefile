.PHONY: build
build:
	sam build

build-MyrestoSettingsFunction:
	GOOS=linux CGO_ENABLED=0 go build -tags lambda.norpc -o $(ARTIFACTS_DIR)/bootstrap .

.PHONY: init
init: build
	sam deploy --guided

.PHONY: deploy
deploy: build
	sam deploy --parameter-overrides \
	ParameterKey=AwsCfToken,ParameterValue="$$AWS_CF_TOKEN" \
	ParameterKey=AdminKey,ParameterValue="$$ADMIN_KEY"

.PHONY: delete
delete:
	sam delete

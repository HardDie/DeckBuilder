.PHONY: default
default: help

.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: swagger
swagger: ## generate swagger json file
	./swagger generate spec -m -o internal/api/web/web/swagger.json

.PHONY: swagger-install-linux
swagger-install-linux: ## install swagger for linux
	curl -o swagger -L https://github.com/go-swagger/go-swagger/releases/download/v0.29.0/swagger_linux_amd64 && chmod +x swagger

.PHONY: linter-install
linter-install: ## install linters
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin

.PHONY: linter-run
linter-run: ## run linters
	./bin/gosec -fmt=sonarqube ./... || echo "gosec found issues"
	./bin/golangci-lint run

.PHONY: test
test: ## run unit tests
	rm -rf data_test || 1
	TEST_DATA_PATH=${PWD}/data_test go test ./... -v -race
	rm -rf data_test || 1

.PHONY: fuzz_game
fuzz_game: ## run fuzzy test for game API's
	rm -rf data_test || 1
	cd internal/service && TEST_DATA_PATH=${PWD}/data_test go test -fuzz=FuzzGame -v
	rm -rf data_test || 1

.PHONY: fuzz_collection
fuzz_collection: ## run fuzzy test for collection API's
	rm -rf data_test || 1
	cd internal/service && TEST_DATA_PATH=${PWD}/data_test go test -fuzz=FuzzCollection -v
	rm -rf data_test || 1

.PHONY: fuzz_deck
fuzz_deck: ## run fuzzy test for deck API's
	rm -rf data_test || 1
	cd internal/service && TEST_DATA_PATH=${PWD}/data_test go test -fuzz=FuzzDeck -v
	rm -rf data_test || 1

.PHONY: fuzz_card
fuzz_card: ## run fuzzy test for card API's
	rm -rf data_test || 1
	cd internal/service && TEST_DATA_PATH=${PWD}/data_test go test -fuzz=FuzzCard -v
	rm -rf data_test || 1

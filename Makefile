run-server:
	go run cmd/tts_deck_builder/main.go

build-windows:
	cd cmd/tts_deck_builder && GOOS=windows go build -o TTS_Deck_Builder.exe .

swagger-spec:
	./swagger generate spec -m -o api/web/web/swagger.json

swagger-lin:
	curl -o swagger -L https://github.com/go-swagger/go-swagger/releases/download/v0.29.0/swagger_linux_amd64 && chmod +x swagger

web-build:
	make -C web

web-setup:
	cp -r web/dist api/web/web

linter-install:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./bin

linter-run: linter-gosec linter-golangci-lint

linter-gosec:
	./bin/gosec -fmt=sonarqube ./... || echo "gosec found issues"

linter-golangci-lint:
	./bin/golangci-lint run

test:
	rm -rf data_test || 1
	TEST_DATA_PATH=${PWD}/data_test go test ./... -v -race
	rm -rf data_test || 1

fuzz_game:
	rm -rf data_test || 1
	cd internal/service && TEST_DATA_PATH=${PWD}/data_test go test -fuzz=FuzzGame -v
	rm -rf data_test || 1

fuzz_collection:
	rm -rf data_test || 1
	cd internal/service && TEST_DATA_PATH=${PWD}/data_test go test -fuzz=FuzzCollection -v
	rm -rf data_test || 1

fuzz_deck:
	rm -rf data_test || 1
	cd internal/service && TEST_DATA_PATH=${PWD}/data_test go test -fuzz=FuzzDeck -v
	rm -rf data_test || 1

fuzz_card:
	rm -rf data_test || 1
	cd internal/service && TEST_DATA_PATH=${PWD}/data_test go test -fuzz=FuzzCard -v
	rm -rf data_test || 1

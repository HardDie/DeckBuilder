run-server:
	go run cmd/tts_deck_builder/main.go

swagger-spec:
	./swagger generate spec -m -o swagger.json
	cp swagger.json api/web/web/swagger.json

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

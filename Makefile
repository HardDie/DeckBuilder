run-server:
	go run cmd/tts_deck_builder/main.go

swagger-spec:
	./swagger generate spec -m -o web/swagger.json

swagger-lin:
	curl -o swagger -L https://github.com/go-swagger/go-swagger/releases/download/v0.29.0/swagger_linux_amd64 && chmod +x swagger

web-build:
	make -C web

web-setup:
	cp -r web/dist api/web/web

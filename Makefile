.PHONY: run build test vet fmt docker-build docker-dev

run:
	go run ./cmd/api

build:
	go build -o bin/ ./cmd/...

test:
	go test ./...

vet:
	go vet ./...

fmt:
	gofmt -w cmd config db internal

docker-build:
	docker build -t social-api .

# Lance l'API en local avec hot-reload (bind mount du code) et l'expose sur le port 8090.
docker-dev:
	docker run --rm -it \
		-p 8090:8080 \
		--env-file .env \
		-e HTTP_ADDR=0.0.0.0:8080 \
		-v "$(CURDIR)":/app \
		-v /app/tmp \
		social-api

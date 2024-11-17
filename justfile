lint:
    go fmt ./...

run:
    APP_ENV=development go run cmd/main.go

build:
    @CGO_ENABLED=0 go build cmd/main.go

test:
    go test ./...

compose:
    docker-compose up -d

down:
    docker-compose down

deploy:
    flyctl deploy
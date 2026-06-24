# Atajos de desarrollo. Uso: make <objetivo>

.PHONY: tidy run test cover docker up down logs

tidy:            ## Resolver dependencias y generar go.sum
	go mod tidy

run:             ## Correr local (SQLite por defecto)
	go run ./cmd/cafeteria-api

test:            ## Correr la suite de tests
	go test ./...

cover:           ## Tests con cobertura
	go test ./... -cover

docker:          ## Construir solo la imagen de la API
	docker build -t cafeteria-uleam-api .

up:              ## Levantar API + PostgreSQL con docker compose
	docker compose up --build

down:            ## Bajar y borrar el volumen de datos
	docker compose down -v

logs:            ## Ver logs de la API
	docker compose logs -f api

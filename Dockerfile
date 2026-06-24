# ============================================================================
# Build multi-stage: una etapa "builder" compila el binario y una etapa final
# minima solo lo copia. Resultado: imagen pequena y sin el toolchain de Go.
# ============================================================================

# ---- Etapa 1: builder ----
FROM golang:1.22-alpine AS builder
WORKDIR /src

# Cachear dependencias: copiar primero los modulos y descargar.
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del codigo y compilar.
COPY . .
# CGO_ENABLED=0 produce un binario estatico (los drivers de SQLite y Postgres
# que usamos son Go puro, asi que no hace falta CGO). GOOS=linux para el runner.
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/cafeteria-api ./cmd/cafeteria-api

# ---- Etapa 2: runner (imagen final minima) ----
FROM alpine:3.20
# ca-certificates por si en el futuro se conecta por TLS; tzdata para zonas horarias.
RUN apk add --no-cache ca-certificates tzdata
# Usuario no-root por seguridad.
RUN adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /bin/cafeteria-api /app/cafeteria-api
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/cafeteria-api"]

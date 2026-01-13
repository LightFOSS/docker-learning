# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar archivos de módulo y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o shift-manager ./cmd/shift-manager

# Stage 2: Run
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copiar el binario compilado desde el stage de build
COPY --from=builder /app/shift-manager .

# Crear el directorio para los datos
RUN mkdir -p /root/.shift-manager

# Volumen para persistir los datos
VOLUME ["/root/.shift-manager"]

# Comando por defecto
ENTRYPOINT ["./shift-manager"]
CMD ["help"]

# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Instalar dependencias necesarias
RUN apk add --no-cache git

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos del módulo Go
COPY go.mod ./

# Descargar dependencias
RUN go mod download

# Copiar el código fuente
COPY main.go ./

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gym-tracker .

# Stage 2: Runtime
FROM alpine:latest

# Instalar certificados CA para HTTPS si es necesario
RUN apk --no-cache add ca-certificates

# Crear usuario no-root
RUN addgroup -g 1000 gymuser && \
    adduser -D -u 1000 -G gymuser gymuser

# Crear directorio para datos
RUN mkdir -p /data && chown gymuser:gymuser /data

# Establecer el directorio de trabajo
WORKDIR /home/gymuser

# Copiar el binario desde el stage de build
COPY --from=builder /app/gym-tracker .

# Cambiar al usuario no-root
USER gymuser

# Establecer variable de entorno para el directorio de datos
ENV DATA_DIR=/data

# Punto de entrada
ENTRYPOINT ["./gym-tracker"]

# Comando por defecto (mostrar ayuda)
CMD []

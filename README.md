# Mi Proyecto Docker

Proyecto de aprendizaje de Docker - Servidor web con Nginx

## ¿Qué hace?
- Servidor web Nginx corriendo en contenedor
- Página HTML personalizada

## Cómo usar
```bash
docker build -t mi-nginx-personalizado:v1 .
docker run -d -p 8080:80 --name nginx-custom mi-nginx-personalizado:v1
```

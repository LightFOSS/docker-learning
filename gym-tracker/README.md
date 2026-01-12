# Gym Tracker CLI

Aplicación de línea de comandos en Go para trackear tus ejercicios de gimnasio, registrar peso, repeticiones y series, y calcular tu 1RM (One Rep Max).

## Características

- ✅ Registrar ejercicios con peso, repeticiones y series
- ✅ Ver historial completo de cualquier ejercicio
- ✅ Calcular 1RM (One Rep Max) automáticamente usando la fórmula de Epley
- ✅ Listar todos los ejercicios registrados
- ✅ Almacenamiento persistente en JSON
- ✅ Containerizado con Docker

## Instalación

### Opción 1: Ejecutar nativamente

#### Requisitos
- Go 1.21 o superior

#### Compilar
```bash
cd gym-tracker
go build -o gym-tracker main.go
```

#### Usar
```bash
./gym-tracker [comando] [opciones]
```

### Opción 2: Ejecutar con Docker (Recomendado)

#### Construir la imagen
```bash
cd gym-tracker
docker build -t gym-tracker .
```

#### Configurar alias (facilita el uso)
```bash
# En Linux/Mac
alias gym='docker run --rm -v gym-data:/data gym-tracker'

# En Windows (PowerShell)
function gym { docker run --rm -v gym-data:/data gym-tracker $args }
```

Ahora puedes usar simplemente: `gym [comando] [opciones]`

## Comandos

### 1. Registrar un ejercicio
```bash
gym-tracker add -name "Press Banca" -weight 80 -reps 5 -sets 3
```

**Parámetros:**
- `-name`: Nombre del ejercicio (requerido)
- `-weight`: Peso levantado en kg (requerido)
- `-reps`: Número de repeticiones (requerido)
- `-sets`: Número de series (requerido)

**Ejemplo con Docker:**
```bash
gym add -name "Sentadilla" -weight 100 -reps 8 -sets 4
```

### 2. Ver historial de un ejercicio
```bash
gym-tracker history -name "Press Banca"
```

Muestra todos los entrenamientos registrados para ese ejercicio, ordenados por fecha (más reciente primero).

**Ejemplo con Docker:**
```bash
gym history -name "Sentadilla"
```

### 3. Calcular 1RM
```bash
gym-tracker 1rm -name "Press Banca"
```

Muestra:
- Mejor 1RM registrado
- Promedio de todos los 1RM calculados
- Detalles del mejor levantamiento

**Ejemplo con Docker:**
```bash
gym 1rm -name "Peso Muerto"
```

### 4. Listar todos los ejercicios
```bash
gym-tracker list
```

Muestra todos los ejercicios únicos registrados con el número de entrenamientos de cada uno.

**Ejemplo con Docker:**
```bash
gym list
```

## Fórmula 1RM

La aplicación utiliza la **fórmula de Epley** para calcular el 1RM:

```
1RM = peso × (1 + repeticiones / 30)
```

Esta fórmula es ampliamente utilizada y proporciona una estimación precisa del peso máximo que podrías levantar en una sola repetición.

## Almacenamiento de datos

### Ejecución nativa
Los datos se guardan en:
- Linux/Mac: `~/.gym-tracker/gym_data.json`
- Windows: `%USERPROFILE%\.gym-tracker\gym_data.json`

### Ejecución con Docker
Los datos se guardan en un volumen de Docker llamado `gym-data`, lo que permite persistencia entre ejecuciones del contenedor.

Para respaldar tus datos:
```bash
# Exportar datos
docker run --rm -v gym-data:/data -v $(pwd):/backup alpine cp /data/gym_data.json /backup/

# Importar datos
docker run --rm -v gym-data:/data -v $(pwd):/backup alpine cp /backup/gym_data.json /data/
```

## Ejemplos de uso completo

```bash
# Registrar una sesión de press banca
gym add -name "Press Banca" -weight 80 -reps 5 -sets 5

# Registrar una sesión de sentadilla
gym add -name "Sentadilla" -weight 100 -reps 8 -sets 4

# Registrar peso muerto
gym add -name "Peso Muerto" -weight 120 -reps 3 -sets 3

# Ver todos los ejercicios
gym list

# Ver historial de press banca
gym history -name "Press Banca"

# Ver mejor 1RM de sentadilla
gym 1rm -name "Sentadilla"
```

## Estructura del proyecto

```
gym-tracker/
├── main.go          # Código fuente de la aplicación
├── go.mod           # Definición del módulo Go
├── Dockerfile       # Configuración de Docker
├── .dockerignore    # Archivos excluidos del build
└── README.md        # Este archivo
```

## Tecnologías utilizadas

- **Go 1.21**: Lenguaje de programación
- **Docker**: Containerización
- **Alpine Linux**: Imagen base ligera
- **JSON**: Formato de almacenamiento de datos

## Ventajas de usar Docker

1. **Sin dependencias locales**: No necesitas instalar Go
2. **Portabilidad**: Funciona igual en cualquier sistema
3. **Aislamiento**: No interfiere con otras aplicaciones
4. **Persistencia**: Los datos se mantienen en un volumen Docker
5. **Ligero**: Imagen final de ~10MB

## Contribuir

Si encuentras algún bug o tienes sugerencias de mejora, no dudes en abrir un issue o enviar un pull request.

## Licencia

MIT License

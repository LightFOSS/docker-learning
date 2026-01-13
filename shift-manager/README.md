# Shift Manager

🗓️ Aplicación CLI en Go para gestionar turnos de trabajo y eventos personales.

## Características

- ✅ Agregar turnos de trabajo con fecha, hora de entrada/salida y tipo (normal/partido)
- ✅ Agregar eventos personales (gym, terapia, clases, etc.)
- ✅ Ver calendario semanal y mensual con vista visual
- ✅ Listar todos los turnos y eventos con filtros
- ✅ Detección automática de conflictos de horarios
- ✅ Calcular total de horas trabajadas por semana/mes
- ✅ Almacenamiento local en JSON
- 🐳 Completamente containerizable con Docker

## Instalación

### Opción 1: Compilar desde el código fuente

```bash
# Clonar el repositorio
git clone https://github.com/LightFOSS/docker-learning.git
cd docker-learning/shift-manager

# Compilar
go build -o shift-manager ./cmd/shift-manager

# Mover a un directorio en el PATH (opcional)
sudo mv shift-manager /usr/local/bin/
```

### Opción 2: Usar Docker

```bash
# Construir la imagen
docker build -t shift-manager:latest .

# Crear un alias para facilitar el uso
alias shift-manager='docker run --rm -v ~/.shift-manager:/root/.shift-manager shift-manager:latest'
```

## Uso

### Comandos disponibles

#### 1. Agregar turno de trabajo

```bash
shift-manager add-shift -date 2026-01-15 -start 09:00 -end 17:00 -type normal
shift-manager add-shift -date 2026-01-16 -start 14:00 -end 22:00 -type partido -notes "Turno tarde"
```

**Opciones:**
- `-date`: Fecha del turno (formato: YYYY-MM-DD) - **obligatorio**
- `-start`: Hora de entrada (formato: HH:MM) - **obligatorio**
- `-end`: Hora de salida (formato: HH:MM) - **obligatorio**
- `-type`: Tipo de turno (`normal` o `partido`) - por defecto: `normal`
- `-notes`: Notas adicionales - opcional

#### 2. Agregar evento personal

```bash
shift-manager add-event -date 2026-01-15 -start 18:00 -end 19:30 -type gym -title "Entrenamiento"
shift-manager add-event -date 2026-01-17 -start 16:00 -end 17:00 -type terapia -title "Sesión terapia" -notes "Dr. Smith"
```

**Opciones:**
- `-date`: Fecha del evento (formato: YYYY-MM-DD) - **obligatorio**
- `-start`: Hora de inicio (formato: HH:MM) - **obligatorio**
- `-end`: Hora de fin (formato: HH:MM) - **obligatorio**
- `-type`: Tipo de evento (gym, terapia, clases, etc.) - **obligatorio**
- `-title`: Título del evento - **obligatorio**
- `-notes`: Notas adicionales - opcional

#### 3. Ver calendario

```bash
# Ver calendario de la semana actual
shift-manager view

# Ver calendario del mes actual
shift-manager view -period month

# Ver calendario de una semana específica
shift-manager view -period week -date 2026-01-15

# Ver calendario de un mes específico
shift-manager view -period month -date 2026-02-01
```

**Opciones:**
- `-period`: Período a mostrar (`week` o `month`) - por defecto: `week`
- `-date`: Fecha de referencia (formato: YYYY-MM-DD) - por defecto: hoy

#### 4. Listar turnos y eventos

```bash
# Listar todo
shift-manager list

# Listar solo turnos
shift-manager list -type shifts

# Listar solo eventos
shift-manager list -type events

# Listar con rango de fechas
shift-manager list -from 2026-01-01 -to 2026-01-31

# Combinar filtros
shift-manager list -type shifts -from 2026-01-15 -to 2026-01-31
```

**Opciones:**
- `-type`: Filtrar por tipo (`all`, `shifts`, `events`) - por defecto: `all`
- `-from`: Fecha desde (formato: YYYY-MM-DD) - opcional
- `-to`: Fecha hasta (formato: YYYY-MM-DD) - opcional

#### 5. Ver estadísticas

```bash
# Estadísticas de la semana actual
shift-manager stats

# Estadísticas del mes actual
shift-manager stats -period month

# Estadísticas totales
shift-manager stats -period all

# Estadísticas de una semana específica
shift-manager stats -period week -date 2026-01-15
```

**Opciones:**
- `-period`: Período de estadísticas (`week`, `month`, `all`) - por defecto: `week`
- `-date`: Fecha de referencia (formato: YYYY-MM-DD) - por defecto: hoy

#### 6. Ayuda y versión

```bash
# Ver ayuda completa
shift-manager help

# Ver versión
shift-manager version
```

## Ejemplos de uso

### Caso de uso completo

```bash
# 1. Agregar turnos de la semana
shift-manager add-shift -date 2026-01-13 -start 09:00 -end 17:00 -type normal
shift-manager add-shift -date 2026-01-14 -start 09:00 -end 17:00 -type normal
shift-manager add-shift -date 2026-01-15 -start 14:00 -end 22:00 -type partido

# 2. Agregar eventos personales
shift-manager add-event -date 2026-01-13 -start 18:00 -end 19:30 -type gym -title "Gimnasio"
shift-manager add-event -date 2026-01-16 -start 16:00 -end 17:00 -type terapia -title "Terapia"

# 3. Ver el calendario de la semana
shift-manager view

# 4. Ver estadísticas
shift-manager stats

# 5. Listar todas las actividades
shift-manager list
```

## Detección de conflictos

La aplicación detecta automáticamente conflictos de horarios cuando:
- Un nuevo turno se superpone con un turno existente
- Un nuevo evento se superpone con un turno existente
- Un nuevo turno se superpone con un evento existente
- Un nuevo evento se superpone con otro evento

Cuando se detecta un conflicto, se muestra una advertencia pero se permite guardar el turno/evento.

## Almacenamiento de datos

Los datos se guardan en formato JSON en:
- **Local**: `~/.shift-manager/data.json`
- **Docker**: Se persisten mediante un volumen montado

## Uso con Docker

### Crear alias persistente

Agregar a tu `~/.bashrc` o `~/.zshrc`:

```bash
alias shift-manager='docker run --rm -v ~/.shift-manager:/root/.shift-manager shift-manager:latest'
```

### Ejemplos con Docker

```bash
# Construir la imagen
docker build -t shift-manager:latest .

# Usar con volumen para persistir datos
docker run --rm -v ~/.shift-manager:/root/.shift-manager shift-manager:latest add-shift -date 2026-01-15 -start 09:00 -end 17:00 -type normal

# Ver ayuda
docker run --rm shift-manager:latest help

# Ver calendario
docker run --rm -v ~/.shift-manager:/root/.shift-manager shift-manager:latest view -period week
```

## Estructura del proyecto

```
shift-manager/
├── cmd/
│   └── shift-manager/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── models/
│   │   └── models.go            # Modelos de datos
│   ├── storage/
│   │   └── storage.go           # Sistema de almacenamiento
│   ├── commands/
│   │   ├── add_shift.go         # Comando add-shift
│   │   ├── add_event.go         # Comando add-event
│   │   ├── view.go              # Comando view
│   │   ├── stats.go             # Comando stats
│   │   └── list.go              # Comando list
│   └── utils/
│       └── utils.go             # Utilidades
├── Dockerfile                   # Dockerfile multi-stage
├── go.mod
├── go.sum
└── README.md
```

## Tecnologías utilizadas

- **Go 1.21+**: Lenguaje de programación
- **Docker**: Containerización
- **JSON**: Almacenamiento de datos

## Contribuir

Las contribuciones son bienvenidas. Por favor:

1. Haz fork del repositorio
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crea un Pull Request

## Licencia

Este proyecto es parte del repositorio de aprendizaje de Docker.

## Autor

LightFOSS

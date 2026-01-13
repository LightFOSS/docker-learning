package main

import (
	"fmt"
	"os"

	"github.com/LightFOSS/docker-learning/shift-manager/internal/commands"
)

const version = "1.0.0"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error

	switch command {
	case "add-shift":
		err = commands.AddShift(args)
	case "add-event":
		err = commands.AddEvent(args)
	case "list":
		err = commands.List(args)
	case "view":
		err = commands.View(args)
	case "stats":
		err = commands.Stats(args)
	case "help", "-h", "--help":
		printHelp()
	case "version", "-v", "--version":
		fmt.Printf("shift-manager v%s\n", version)
	default:
		fmt.Printf("Comando desconocido: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Uso: shift-manager <comando> [opciones]")
	fmt.Println()
	fmt.Println("Comandos disponibles:")
	fmt.Println("  add-shift    Agregar un turno de trabajo")
	fmt.Println("  add-event    Agregar un evento personal")
	fmt.Println("  list         Listar todos los turnos y eventos")
	fmt.Println("  view         Ver calendario semanal o mensual")
	fmt.Println("  stats        Ver estadísticas de horas trabajadas")
	fmt.Println("  help         Mostrar ayuda detallada")
	fmt.Println("  version      Mostrar versión")
	fmt.Println()
	fmt.Println("Usa 'shift-manager help' para ver más información sobre los comandos.")
}

func printHelp() {
	fmt.Println("shift-manager - Gestión de turnos de trabajo y eventos personales")
	fmt.Println()
	fmt.Println("COMANDOS:")
	fmt.Println()

	fmt.Println("  add-shift - Agregar un turno de trabajo")
	fmt.Println("    Uso: shift-manager add-shift -date YYYY-MM-DD -start HH:MM -end HH:MM -type TYPE [-notes NOTES]")
	fmt.Println("    Opciones:")
	fmt.Println("      -date    Fecha del turno (obligatorio)")
	fmt.Println("      -start   Hora de entrada (obligatorio)")
	fmt.Println("      -end     Hora de salida (obligatorio)")
	fmt.Println("      -type    Tipo de turno: normal o partido (por defecto: normal)")
	fmt.Println("      -notes   Notas adicionales (opcional)")
	fmt.Println()
	fmt.Println("    Ejemplo:")
	fmt.Println("      shift-manager add-shift -date 2026-01-15 -start 09:00 -end 17:00 -type normal")
	fmt.Println()

	fmt.Println("  add-event - Agregar un evento personal")
	fmt.Println("    Uso: shift-manager add-event -date YYYY-MM-DD -start HH:MM -end HH:MM -type TYPE -title TITLE [-notes NOTES]")
	fmt.Println("    Opciones:")
	fmt.Println("      -date    Fecha del evento (obligatorio)")
	fmt.Println("      -start   Hora de inicio (obligatorio)")
	fmt.Println("      -end     Hora de fin (obligatorio)")
	fmt.Println("      -type    Tipo de evento: gym, terapia, clases, etc. (obligatorio)")
	fmt.Println("      -title   Título del evento (obligatorio)")
	fmt.Println("      -notes   Notas adicionales (opcional)")
	fmt.Println()
	fmt.Println("    Ejemplo:")
	fmt.Println("      shift-manager add-event -date 2026-01-15 -start 18:00 -end 19:30 -type gym -title 'Entrenamiento'")
	fmt.Println()

	fmt.Println("  list - Listar todos los turnos y eventos")
	fmt.Println("    Uso: shift-manager list [-type TYPE] [-from YYYY-MM-DD] [-to YYYY-MM-DD]")
	fmt.Println("    Opciones:")
	fmt.Println("      -type    Filtrar por tipo: all, shifts, events (por defecto: all)")
	fmt.Println("      -from    Fecha desde (opcional)")
	fmt.Println("      -to      Fecha hasta (opcional)")
	fmt.Println()
	fmt.Println("    Ejemplos:")
	fmt.Println("      shift-manager list")
	fmt.Println("      shift-manager list -type shifts")
	fmt.Println("      shift-manager list -from 2026-01-01 -to 2026-01-31")
	fmt.Println()

	fmt.Println("  view - Ver calendario semanal o mensual")
	fmt.Println("    Uso: shift-manager view [-period PERIOD] [-date YYYY-MM-DD]")
	fmt.Println("    Opciones:")
	fmt.Println("      -period  Período: week o month (por defecto: week)")
	fmt.Println("      -date    Fecha de referencia (por defecto: hoy)")
	fmt.Println()
	fmt.Println("    Ejemplos:")
	fmt.Println("      shift-manager view")
	fmt.Println("      shift-manager view -period month")
	fmt.Println("      shift-manager view -period week -date 2026-01-15")
	fmt.Println()

	fmt.Println("  stats - Ver estadísticas de horas trabajadas")
	fmt.Println("    Uso: shift-manager stats [-period PERIOD] [-date YYYY-MM-DD]")
	fmt.Println("    Opciones:")
	fmt.Println("      -period  Período: week, month o all (por defecto: week)")
	fmt.Println("      -date    Fecha de referencia (por defecto: hoy)")
	fmt.Println()
	fmt.Println("    Ejemplos:")
	fmt.Println("      shift-manager stats")
	fmt.Println("      shift-manager stats -period month")
	fmt.Println("      shift-manager stats -period all")
	fmt.Println()

	fmt.Println("CARACTERÍSTICAS:")
	fmt.Println("  • Detección automática de conflictos de horarios")
	fmt.Println("  • Cálculo de horas trabajadas por turno, semana y mes")
	fmt.Println("  • Calendario visual semanal y mensual")
	fmt.Println("  • Soporte para turnos normales y partidos")
	fmt.Println("  • Almacenamiento local en JSON")
	fmt.Println()
	fmt.Println("Los datos se guardan en: ~/.shift-manager/data.json")
}

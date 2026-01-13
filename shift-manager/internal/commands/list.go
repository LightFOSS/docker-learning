package commands

import (
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/LightFOSS/docker-learning/shift-manager/internal/storage"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/utils"
)

type listItem struct {
	Date      time.Time
	Type      string
	StartTime string
	EndTime   string
	Details   string
}

// List lista todos los turnos y eventos
func List(args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	filterType := fs.String("type", "all", "Filtrar por tipo (all/shifts/events)")
	fromDate := fs.String("from", "", "Fecha desde (YYYY-MM-DD)")
	toDate := fs.String("to", "", "Fecha hasta (YYYY-MM-DD)")

	if err := fs.Parse(args); err != nil {
		return err
	}

	store, err := storage.New()
	if err != nil {
		return err
	}

	data, err := store.Load()
	if err != nil {
		return err
	}

	// Validar fechas si se proporcionan
	var from, to time.Time
	if *fromDate != "" {
		from, err = utils.ParseDate(*fromDate)
		if err != nil {
			return fmt.Errorf("fecha 'from' inválida: %w", err)
		}
	}

	if *toDate != "" {
		to, err = utils.ParseDate(*toDate)
		if err != nil {
			return fmt.Errorf("fecha 'to' inválida: %w", err)
		}
	}

	// Crear lista combinada
	var items []listItem

	// Agregar turnos
	if *filterType == "all" || *filterType == "shifts" {
		for _, shift := range data.Shifts {
			date, err := utils.ParseDate(shift.Date)
			if err != nil {
				continue
			}

			// Filtrar por rango de fechas
			if *fromDate != "" && date.Before(from) {
				continue
			}
			if *toDate != "" && date.After(to) {
				continue
			}

			hours, _ := shift.CalculateHours()
			details := fmt.Sprintf("Turno %s - %s", shift.Type, utils.FormatDuration(hours))
			if shift.Notes != "" {
				details += fmt.Sprintf(" | %s", shift.Notes)
			}

			items = append(items, listItem{
				Date:      date,
				Type:      "TURNO",
				StartTime: shift.StartTime,
				EndTime:   shift.EndTime,
				Details:   details,
			})
		}
	}

	// Agregar eventos
	if *filterType == "all" || *filterType == "events" {
		for _, event := range data.Events {
			date, err := utils.ParseDate(event.Date)
			if err != nil {
				continue
			}

			// Filtrar por rango de fechas
			if *fromDate != "" && date.Before(from) {
				continue
			}
			if *toDate != "" && date.After(to) {
				continue
			}

			details := fmt.Sprintf("%s (%s)", event.Title, event.Type)
			if event.Notes != "" {
				details += fmt.Sprintf(" | %s", event.Notes)
			}

			items = append(items, listItem{
				Date:      date,
				Type:      "EVENTO",
				StartTime: event.StartTime,
				EndTime:   event.EndTime,
				Details:   details,
			})
		}
	}

	// Ordenar por fecha
	sort.Slice(items, func(i, j int) bool {
		if items[i].Date.Equal(items[j].Date) {
			return items[i].StartTime < items[j].StartTime
		}
		return items[i].Date.Before(items[j].Date)
	})

	// Mostrar resultados
	if len(items) == 0 {
		fmt.Println("No se encontraron turnos ni eventos.")
		return nil
	}

	fmt.Printf("\n📋 Lista de %s\n", getFilterDescription(*filterType))
	if *fromDate != "" || *toDate != "" {
		fmt.Printf("   Período: %s - %s\n", getDateOrDefault(*fromDate, "inicio"), getDateOrDefault(*toDate, "fin"))
	}
	fmt.Println()

	currentDate := ""
	for _, item := range items {
		dateStr := utils.FormatDate(item.Date)
		if dateStr != currentDate {
			if currentDate != "" {
				fmt.Println()
			}
			fmt.Printf("📅 %s (%s)\n", dateStr, item.Date.Weekday().String())
			currentDate = dateStr
		}

		fmt.Printf("   [%s] %s - %s | %s\n", item.Type, item.StartTime, item.EndTime, item.Details)
	}

	fmt.Printf("\n📊 Total: %d item(s)\n", len(items))

	return nil
}

func getFilterDescription(filterType string) string {
	switch filterType {
	case "shifts":
		return "Turnos"
	case "events":
		return "Eventos"
	default:
		return "Turnos y Eventos"
	}
}

func getDateOrDefault(date, defaultVal string) string {
	if date == "" {
		return defaultVal
	}
	return date
}

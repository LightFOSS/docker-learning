package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/LightFOSS/docker-learning/shift-manager/internal/models"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/storage"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/utils"
)

// View muestra el calendario de turnos y eventos
func View(args []string) error {
	fs := flag.NewFlagSet("view", flag.ExitOnError)
	period := fs.String("period", "week", "Período a mostrar (week/month)")
	date := fs.String("date", "", "Fecha de referencia (YYYY-MM-DD), por defecto hoy")

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

	// Determinar fecha de referencia
	var refDate time.Time
	if *date != "" {
		refDate, err = utils.ParseDate(*date)
		if err != nil {
			return fmt.Errorf("fecha inválida: %w", err)
		}
	} else {
		refDate = time.Now()
	}

	// Mostrar calendario según el período
	switch *period {
	case "week":
		return showWeekView(data, refDate)
	case "month":
		return showMonthView(data, refDate)
	default:
		return fmt.Errorf("período inválido (debe ser week o month)")
	}
}

func showWeekView(data *models.Data, refDate time.Time) error {
	startOfWeek, endOfWeek := utils.GetWeekBounds(refDate)

	fmt.Printf("\n📅 Calendario Semanal\n")
	fmt.Printf("   Semana del %s al %s\n\n", utils.FormatDate(startOfWeek), utils.FormatDate(endOfWeek))

	// Mostrar cada día de la semana
	for d := 0; d < 7; d++ {
		currentDay := startOfWeek.AddDate(0, 0, d)
		dayStr := utils.FormatDate(currentDay)

		// Marcar el día actual
		marker := "  "
		if utils.FormatDate(currentDay) == utils.FormatDate(time.Now()) {
			marker = "→ "
		}

		fmt.Printf("%s%s (%s)\n", marker, dayStr, currentDay.Weekday().String())

		// Buscar turnos de este día
		hasItems := false
		for _, shift := range data.Shifts {
			if shift.Date == dayStr {
				hours, _ := shift.CalculateHours()
				fmt.Printf("   🏢 Turno %s: %s - %s (%s)\n",
					shift.Type, shift.StartTime, shift.EndTime, utils.FormatDuration(hours))
				if shift.Notes != "" {
					fmt.Printf("      💬 %s\n", shift.Notes)
				}
				hasItems = true
			}
		}

		// Buscar eventos de este día
		for _, event := range data.Events {
			if event.Date == dayStr {
				fmt.Printf("   📌 %s: %s - %s [%s]\n",
					event.Title, event.StartTime, event.EndTime, event.Type)
				if event.Notes != "" {
					fmt.Printf("      💬 %s\n", event.Notes)
				}
				hasItems = true
			}
		}

		if !hasItems {
			fmt.Println("   (sin actividades)")
		}

		fmt.Println()
	}

	return nil
}

func showMonthView(data *models.Data, refDate time.Time) error {
	startOfMonth, endOfMonth := utils.GetMonthBounds(refDate)

	fmt.Printf("\n📅 Calendario Mensual\n")
	fmt.Printf("   %s %d\n\n", refDate.Month().String(), refDate.Year())

	// Encabezado del calendario
	fmt.Println("   Lun  Mar  Mié  Jue  Vie  Sáb  Dom")
	fmt.Println("   ─────────────────────────────────")

	// Determinar el primer día del mes y su día de la semana
	firstDay := startOfMonth
	weekday := int(firstDay.Weekday())
	if weekday == 0 {
		weekday = 7 // Domingo al final
	}

	// Espacios iniciales
	fmt.Print("  ")
	for i := 1; i < weekday; i++ {
		fmt.Print("     ")
	}

	// Días del mes
	currentWeekday := weekday
	for day := 1; day <= endOfMonth.Day(); day++ {
		currentDate := time.Date(refDate.Year(), refDate.Month(), day, 0, 0, 0, 0, refDate.Location())
		dateStr := utils.FormatDate(currentDate)

		// Verificar si hay actividades este día
		hasShift := false
		hasEvent := false

		for _, shift := range data.Shifts {
			if shift.Date == dateStr {
				hasShift = true
				break
			}
		}

		for _, event := range data.Events {
			if event.Date == dateStr {
				hasEvent = true
				break
			}
		}

		// Determinar el símbolo del día
		symbol := fmt.Sprintf("%2d", day)
		if utils.FormatDate(currentDate) == utils.FormatDate(time.Now()) {
			symbol = fmt.Sprintf("[%2d]", day)
		} else if hasShift && hasEvent {
			symbol = fmt.Sprintf("*%2d*", day)
		} else if hasShift {
			symbol = fmt.Sprintf("·%2d·", day)
		} else if hasEvent {
			symbol = fmt.Sprintf("°%2d°", day)
		} else {
			symbol = fmt.Sprintf(" %2d ", day)
		}

		fmt.Printf(" %s", symbol)

		currentWeekday++
		if currentWeekday > 7 {
			fmt.Println()
			fmt.Print("  ")
			currentWeekday = 1
		}
	}
	fmt.Println("\n")

	// Leyenda
	fmt.Println("Leyenda:")
	fmt.Println("  [XX] = Hoy")
	fmt.Println("  ·XX· = Tiene turno(s)")
	fmt.Println("  °XX° = Tiene evento(s)")
	fmt.Println("  *XX* = Tiene ambos")
	fmt.Println()

	// Mostrar lista de actividades del mes
	fmt.Println("📋 Actividades del mes:")
	hasActivities := false

	// Agrupar por día
	for day := 1; day <= endOfMonth.Day(); day++ {
		currentDate := time.Date(refDate.Year(), refDate.Month(), day, 0, 0, 0, 0, refDate.Location())
		dateStr := utils.FormatDate(currentDate)

		dayActivities := []string{}

		for _, shift := range data.Shifts {
			if shift.Date == dateStr {
				hours, _ := shift.CalculateHours()
				dayActivities = append(dayActivities,
					fmt.Sprintf("    🏢 Turno %s: %s - %s (%s)",
						shift.Type, shift.StartTime, shift.EndTime, utils.FormatDuration(hours)))
			}
		}

		for _, event := range data.Events {
			if event.Date == dateStr {
				dayActivities = append(dayActivities,
					fmt.Sprintf("    📌 %s: %s - %s [%s]",
						event.Title, event.StartTime, event.EndTime, event.Type))
			}
		}

		if len(dayActivities) > 0 {
			fmt.Printf("\n  %s (%s):\n", dateStr, currentDate.Weekday().String()[:3])
			for _, activity := range dayActivities {
				fmt.Println(activity)
			}
			hasActivities = true
		}
	}

	if !hasActivities {
		fmt.Println("  (sin actividades este mes)")
	}

	fmt.Println()

	return nil
}

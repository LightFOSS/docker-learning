package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/LightFOSS/docker-learning/shift-manager/internal/models"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/storage"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/utils"
)

// Stats muestra estadísticas de horas trabajadas
func Stats(args []string) error {
	fs := flag.NewFlagSet("stats", flag.ExitOnError)
	period := fs.String("period", "week", "Período de estadísticas (week/month/all)")
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

	// Calcular estadísticas según el período
	switch *period {
	case "week":
		return showWeekStats(data, refDate)
	case "month":
		return showMonthStats(data, refDate)
	case "all":
		return showAllStats(data)
	default:
		return fmt.Errorf("período inválido (debe ser week, month o all)")
	}
}

func showWeekStats(data *models.Data, refDate time.Time) error {
	startOfWeek, endOfWeek := utils.GetWeekBounds(refDate)

	fmt.Printf("\n📊 Estadísticas Semanales\n")
	fmt.Printf("   Semana del %s al %s\n\n", utils.FormatDate(startOfWeek), utils.FormatDate(endOfWeek))

	var totalHours float64
	var normalHours float64
	var partidoHours float64
	var shiftCount int
	var normalCount int
	var partidoCount int

	for _, shift := range data.Shifts {
		shiftDate, err := utils.ParseDate(shift.Date)
		if err != nil {
			continue
		}

		if shiftDate.Before(startOfWeek) || shiftDate.After(endOfWeek) {
			continue
		}

		hours, err := shift.CalculateHours()
		if err != nil {
			continue
		}

		totalHours += hours
		shiftCount++

		if shift.Type == models.ShiftTypeNormal {
			normalHours += hours
			normalCount++
		} else {
			partidoHours += hours
			partidoCount++
		}
	}

	fmt.Printf("Total de turnos: %d\n", shiftCount)
	fmt.Printf("  • Normales: %d (%s)\n", normalCount, utils.FormatDuration(normalHours))
	fmt.Printf("  • Partidos: %d (%s)\n", partidoCount, utils.FormatDuration(partidoHours))
	fmt.Printf("\nTotal de horas: %s\n", utils.FormatDuration(totalHours))
	fmt.Printf("Promedio diario: %s\n", utils.FormatDuration(totalHours/7))

	// Mostrar desglose por día
	fmt.Println("\n📅 Desglose por día:")
	for d := 0; d < 7; d++ {
		currentDay := startOfWeek.AddDate(0, 0, d)
		dayStr := utils.FormatDate(currentDay)
		dayHours := 0.0
		dayShifts := 0

		for _, shift := range data.Shifts {
			if shift.Date == dayStr {
				hours, _ := shift.CalculateHours()
				dayHours += hours
				dayShifts++
			}
		}

		if dayShifts > 0 {
			fmt.Printf("  %s (%s): %s (%d turno(s))\n",
				dayStr,
				currentDay.Weekday().String()[:3],
				utils.FormatDuration(dayHours),
				dayShifts)
		}
	}

	return nil
}

func showMonthStats(data *models.Data, refDate time.Time) error {
	startOfMonth, endOfMonth := utils.GetMonthBounds(refDate)

	fmt.Printf("\n📊 Estadísticas Mensuales\n")
	fmt.Printf("   %s %d\n\n", refDate.Month().String(), refDate.Year())

	var totalHours float64
	var normalHours float64
	var partidoHours float64
	var shiftCount int
	var normalCount int
	var partidoCount int

	weekStats := make(map[int]float64)

	for _, shift := range data.Shifts {
		shiftDate, err := utils.ParseDate(shift.Date)
		if err != nil {
			continue
		}

		if shiftDate.Before(startOfMonth) || shiftDate.After(endOfMonth) {
			continue
		}

		hours, err := shift.CalculateHours()
		if err != nil {
			continue
		}

		totalHours += hours
		shiftCount++

		if shift.Type == models.ShiftTypeNormal {
			normalHours += hours
			normalCount++
		} else {
			partidoHours += hours
			partidoCount++
		}

		// Agrupar por semana
		_, week := shiftDate.ISOWeek()
		weekStats[week] += hours
	}

	fmt.Printf("Total de turnos: %d\n", shiftCount)
	fmt.Printf("  • Normales: %d (%s)\n", normalCount, utils.FormatDuration(normalHours))
	fmt.Printf("  • Partidos: %d (%s)\n", partidoCount, utils.FormatDuration(partidoHours))
	fmt.Printf("\nTotal de horas: %s\n", utils.FormatDuration(totalHours))

	daysInMonth := endOfMonth.Day()
	fmt.Printf("Promedio diario: %s\n", utils.FormatDuration(totalHours/float64(daysInMonth)))

	// Mostrar desglose por semana
	if len(weekStats) > 0 {
		fmt.Println("\n📅 Desglose por semana:")
		for week, hours := range weekStats {
			fmt.Printf("  Semana %d: %s\n", week, utils.FormatDuration(hours))
		}
	}

	return nil
}

func showAllStats(data *models.Data) error {
	fmt.Printf("\n📊 Estadísticas Totales\n\n")

	if len(data.Shifts) == 0 {
		fmt.Println("No hay turnos registrados.")
		return nil
	}

	var totalHours float64
	var normalHours float64
	var partidoHours float64
	var normalCount int
	var partidoCount int

	var firstDate, lastDate time.Time

	for i, shift := range data.Shifts {
		hours, err := shift.CalculateHours()
		if err != nil {
			continue
		}

		totalHours += hours

		if shift.Type == models.ShiftTypeNormal {
			normalHours += hours
			normalCount++
		} else {
			partidoHours += hours
			partidoCount++
		}

		shiftDate, err := utils.ParseDate(shift.Date)
		if err != nil {
			continue
		}

		if i == 0 {
			firstDate = shiftDate
			lastDate = shiftDate
		} else {
			if shiftDate.Before(firstDate) {
				firstDate = shiftDate
			}
			if shiftDate.After(lastDate) {
				lastDate = shiftDate
			}
		}
	}

	fmt.Printf("Total de turnos: %d\n", len(data.Shifts))
	fmt.Printf("  • Normales: %d (%s)\n", normalCount, utils.FormatDuration(normalHours))
	fmt.Printf("  • Partidos: %d (%s)\n", partidoCount, utils.FormatDuration(partidoHours))
	fmt.Printf("\nTotal de horas: %s\n", utils.FormatDuration(totalHours))

	if !firstDate.IsZero() && !lastDate.IsZero() {
		fmt.Printf("\nPeríodo: %s - %s\n", utils.FormatDate(firstDate), utils.FormatDate(lastDate))
		days := lastDate.Sub(firstDate).Hours()/24 + 1
		if days > 0 {
			fmt.Printf("Promedio diario: %s\n", utils.FormatDuration(totalHours/days))
		}
	}

	fmt.Printf("\nTotal de eventos: %d\n", len(data.Events))

	return nil
}

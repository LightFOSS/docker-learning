package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/LightFOSS/docker-learning/shift-manager/internal/models"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/storage"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/utils"
)

// AddShift agrega un nuevo turno de trabajo
func AddShift(args []string) error {
	fs := flag.NewFlagSet("add-shift", flag.ExitOnError)
	date := fs.String("date", "", "Fecha del turno (YYYY-MM-DD)")
	startTime := fs.String("start", "", "Hora de entrada (HH:MM)")
	endTime := fs.String("end", "", "Hora de salida (HH:MM)")
	shiftType := fs.String("type", "normal", "Tipo de turno (normal/partido)")
	notes := fs.String("notes", "", "Notas adicionales")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validaciones
	if *date == "" || *startTime == "" || *endTime == "" {
		return fmt.Errorf("los parámetros date, start y end son obligatorios")
	}

	if err := utils.ValidateDate(*date); err != nil {
		return err
	}

	if err := utils.ValidateTime(*startTime); err != nil {
		return err
	}

	if err := utils.ValidateTime(*endTime); err != nil {
		return err
	}

	// Validar tipo de turno
	var st models.ShiftType
	switch *shiftType {
	case "normal":
		st = models.ShiftTypeNormal
	case "partido":
		st = models.ShiftTypePartido
	default:
		return fmt.Errorf("tipo de turno inválido (debe ser 'normal' o 'partido')")
	}

	// Crear turno
	shift := models.Shift{
		ID:        utils.GenerateID(),
		Date:      *date,
		StartTime: *startTime,
		EndTime:   *endTime,
		Type:      st,
		Notes:     *notes,
		CreatedAt: time.Now(),
	}

	// Verificar conflictos
	store, err := storage.New()
	if err != nil {
		return err
	}

	data, err := store.Load()
	if err != nil {
		return err
	}

	conflicts, err := checkShiftConflicts(shift, data)
	if err != nil {
		return err
	}

	if len(conflicts) > 0 {
		fmt.Println("⚠️  ADVERTENCIA: Se detectaron conflictos de horario:")
		for _, conflict := range conflicts {
			fmt.Printf("  - %s\n", conflict)
		}
		fmt.Println()
	}

	// Guardar turno
	if err := store.AddShift(shift); err != nil {
		return err
	}

	hours, _ := shift.CalculateHours()
	fmt.Printf("✅ Turno agregado exitosamente!\n")
	fmt.Printf("   ID: %s\n", shift.ID)
	fmt.Printf("   Fecha: %s\n", shift.Date)
	fmt.Printf("   Horario: %s - %s\n", shift.StartTime, shift.EndTime)
	fmt.Printf("   Tipo: %s\n", shift.Type)
	fmt.Printf("   Horas: %s\n", utils.FormatDuration(hours))
	if shift.Notes != "" {
		fmt.Printf("   Notas: %s\n", shift.Notes)
	}

	return nil
}

// checkShiftConflicts verifica si hay conflictos de horario
func checkShiftConflicts(newShift models.Shift, data *models.Data) ([]string, error) {
	var conflicts []string

	newStart, err := newShift.GetStartDateTime()
	if err != nil {
		return nil, err
	}

	newEnd, err := newShift.GetEndDateTime()
	if err != nil {
		return nil, err
	}

	// Verificar conflictos con otros turnos
	for _, shift := range data.Shifts {
		if shift.Date != newShift.Date {
			continue
		}

		start, err := shift.GetStartDateTime()
		if err != nil {
			continue
		}

		end, err := shift.GetEndDateTime()
		if err != nil {
			continue
		}

		if models.HasConflict(newStart, newEnd, start, end) {
			conflicts = append(conflicts, fmt.Sprintf("Turno el %s de %s a %s", shift.Date, shift.StartTime, shift.EndTime))
		}
	}

	// Verificar conflictos con eventos
	for _, event := range data.Events {
		if event.Date != newShift.Date {
			continue
		}

		start, err := event.GetStartDateTime()
		if err != nil {
			continue
		}

		end, err := event.GetEndDateTime()
		if err != nil {
			continue
		}

		if models.HasConflict(newStart, newEnd, start, end) {
			conflicts = append(conflicts, fmt.Sprintf("Evento '%s' el %s de %s a %s", event.Title, event.Date, event.StartTime, event.EndTime))
		}
	}

	return conflicts, nil
}

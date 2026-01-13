package commands

import (
	"flag"
	"fmt"
	"time"

	"github.com/LightFOSS/docker-learning/shift-manager/internal/models"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/storage"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/utils"
)

// AddEvent agrega un nuevo evento personal
func AddEvent(args []string) error {
	fs := flag.NewFlagSet("add-event", flag.ExitOnError)
	date := fs.String("date", "", "Fecha del evento (YYYY-MM-DD)")
	startTime := fs.String("start", "", "Hora de inicio (HH:MM)")
	endTime := fs.String("end", "", "Hora de fin (HH:MM)")
	eventType := fs.String("type", "", "Tipo de evento (gym, terapia, clases, etc.)")
	title := fs.String("title", "", "Título del evento")
	notes := fs.String("notes", "", "Notas adicionales")

	if err := fs.Parse(args); err != nil {
		return err
	}

	// Validaciones
	if *date == "" || *startTime == "" || *endTime == "" || *eventType == "" || *title == "" {
		return fmt.Errorf("los parámetros date, start, end, type y title son obligatorios")
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

	// Crear evento
	event := models.Event{
		ID:        utils.GenerateID(),
		Date:      *date,
		StartTime: *startTime,
		EndTime:   *endTime,
		Type:      *eventType,
		Title:     *title,
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

	conflicts, err := checkEventConflicts(event, data)
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

	// Guardar evento
	if err := store.AddEvent(event); err != nil {
		return err
	}

	fmt.Printf("✅ Evento agregado exitosamente!\n")
	fmt.Printf("   ID: %s\n", event.ID)
	fmt.Printf("   Título: %s\n", event.Title)
	fmt.Printf("   Fecha: %s\n", event.Date)
	fmt.Printf("   Horario: %s - %s\n", event.StartTime, event.EndTime)
	fmt.Printf("   Tipo: %s\n", event.Type)
	if event.Notes != "" {
		fmt.Printf("   Notas: %s\n", event.Notes)
	}

	return nil
}

// checkEventConflicts verifica si hay conflictos de horario
func checkEventConflicts(newEvent models.Event, data *models.Data) ([]string, error) {
	var conflicts []string

	newStart, err := newEvent.GetStartDateTime()
	if err != nil {
		return nil, err
	}

	newEnd, err := newEvent.GetEndDateTime()
	if err != nil {
		return nil, err
	}

	// Verificar conflictos con turnos
	for _, shift := range data.Shifts {
		if shift.Date != newEvent.Date {
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

	// Verificar conflictos con otros eventos
	for _, event := range data.Events {
		if event.Date != newEvent.Date {
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

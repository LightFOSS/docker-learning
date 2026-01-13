package models

import (
	"fmt"
	"time"
)

// ShiftType representa el tipo de turno de trabajo
type ShiftType string

const (
	ShiftTypeNormal  ShiftType = "normal"
	ShiftTypePartido ShiftType = "partido"
)

// Shift representa un turno de trabajo
type Shift struct {
	ID        string    `json:"id"`
	Date      string    `json:"date"`       // Formato: YYYY-MM-DD
	StartTime string    `json:"start_time"` // Formato: HH:MM
	EndTime   string    `json:"end_time"`   // Formato: HH:MM
	Type      ShiftType `json:"type"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Event representa un evento personal
type Event struct {
	ID        string    `json:"id"`
	Date      string    `json:"date"`       // Formato: YYYY-MM-DD
	StartTime string    `json:"start_time"` // Formato: HH:MM
	EndTime   string    `json:"end_time"`   // Formato: HH:MM
	Type      string    `json:"type"`       // gym, terapia, clases, etc.
	Title     string    `json:"title"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// Data contiene todos los turnos y eventos
type Data struct {
	Shifts []Shift `json:"shifts"`
	Events []Event `json:"events"`
}

// GetDateTime combina la fecha y hora en un time.Time
func GetDateTime(date, timeStr string) (time.Time, error) {
	dateTimeStr := fmt.Sprintf("%s %s", date, timeStr)
	return time.Parse("2006-01-02 15:04", dateTimeStr)
}

// CalculateHours calcula las horas trabajadas en un turno
func (s *Shift) CalculateHours() (float64, error) {
	start, err := GetDateTime(s.Date, s.StartTime)
	if err != nil {
		return 0, err
	}

	end, err := GetDateTime(s.Date, s.EndTime)
	if err != nil {
		return 0, err
	}

	// Si el turno termina después de medianoche
	if end.Before(start) {
		end = end.Add(24 * time.Hour)
	}

	duration := end.Sub(start)
	return duration.Hours(), nil
}

// GetStartDateTime retorna el datetime de inicio del turno
func (s *Shift) GetStartDateTime() (time.Time, error) {
	return GetDateTime(s.Date, s.StartTime)
}

// GetEndDateTime retorna el datetime de fin del turno
func (s *Shift) GetEndDateTime() (time.Time, error) {
	end, err := GetDateTime(s.Date, s.EndTime)
	if err != nil {
		return time.Time{}, err
	}

	start, err := GetDateTime(s.Date, s.StartTime)
	if err != nil {
		return time.Time{}, err
	}

	// Si el turno termina después de medianoche
	if end.Before(start) {
		end = end.Add(24 * time.Hour)
	}

	return end, nil
}

// GetStartDateTime retorna el datetime de inicio del evento
func (e *Event) GetStartDateTime() (time.Time, error) {
	return GetDateTime(e.Date, e.StartTime)
}

// GetEndDateTime retorna el datetime de fin del evento
func (e *Event) GetEndDateTime() (time.Time, error) {
	return GetDateTime(e.Date, e.EndTime)
}

// HasConflict verifica si dos rangos de tiempo se superponen
func HasConflict(start1, end1, start2, end2 time.Time) bool {
	// Dos rangos se superponen si:
	// start1 < end2 && start2 < end1
	return start1.Before(end2) && start2.Before(end1)
}

package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// GenerateID genera un ID único
func GenerateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// ValidateDate valida que una fecha tenga el formato correcto
func ValidateDate(date string) error {
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		return fmt.Errorf("formato de fecha inválido (debe ser YYYY-MM-DD): %w", err)
	}
	return nil
}

// ValidateTime valida que una hora tenga el formato correcto
func ValidateTime(timeStr string) error {
	_, err := time.Parse("15:04", timeStr)
	if err != nil {
		return fmt.Errorf("formato de hora inválido (debe ser HH:MM): %w", err)
	}
	return nil
}

// FormatDate formatea un time.Time a string YYYY-MM-DD
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatTime formatea un time.Time a string HH:MM
func FormatTime(t time.Time) string {
	return t.Format("15:04")
}

// ParseDate parsea una fecha en formato YYYY-MM-DD
func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}

// GetWeekBounds retorna el inicio y fin de la semana para una fecha dada
func GetWeekBounds(date time.Time) (time.Time, time.Time) {
	// Lunes como primer día de la semana
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // Domingo
	}

	startOfWeek := date.AddDate(0, 0, -(weekday - 1))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	endOfWeek := startOfWeek.AddDate(0, 0, 6)
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, endOfWeek.Location())

	return startOfWeek, endOfWeek
}

// GetMonthBounds retorna el inicio y fin del mes para una fecha dada
func GetMonthBounds(date time.Time) (time.Time, time.Time) {
	year, month, _ := date.Date()
	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, date.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)
	endOfMonth = time.Date(endOfMonth.Year(), endOfMonth.Month(), endOfMonth.Day(), 23, 59, 59, 0, endOfMonth.Location())

	return startOfMonth, endOfMonth
}

// FormatDuration formatea una duración en horas a formato legible
func FormatDuration(hours float64) string {
	h := int(hours)
	m := int((hours - float64(h)) * 60)
	return fmt.Sprintf("%dh %dm", h, m)
}

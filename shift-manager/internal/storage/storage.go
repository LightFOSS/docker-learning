package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LightFOSS/docker-learning/shift-manager/internal/models"
)

const (
	defaultDataDir  = ".shift-manager"
	defaultDataFile = "data.json"
)

// Storage maneja la persistencia de datos
type Storage struct {
	dataPath string
}

// New crea una nueva instancia de Storage
func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo directorio home: %w", err)
	}

	dataDir := filepath.Join(homeDir, defaultDataDir)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("error creando directorio de datos: %w", err)
	}

	dataPath := filepath.Join(dataDir, defaultDataFile)

	return &Storage{
		dataPath: dataPath,
	}, nil
}

// Load carga los datos desde el archivo JSON
func (s *Storage) Load() (*models.Data, error) {
	data := &models.Data{
		Shifts: []models.Shift{},
		Events: []models.Event{},
	}

	// Si el archivo no existe, retornamos datos vacíos
	if _, err := os.Stat(s.dataPath); os.IsNotExist(err) {
		return data, nil
	}

	file, err := os.ReadFile(s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo archivo: %w", err)
	}

	if len(file) == 0 {
		return data, nil
	}

	if err := json.Unmarshal(file, data); err != nil {
		return nil, fmt.Errorf("error deserializando JSON: %w", err)
	}

	return data, nil
}

// Save guarda los datos en el archivo JSON
func (s *Storage) Save(data *models.Data) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando JSON: %w", err)
	}

	if err := os.WriteFile(s.dataPath, jsonData, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	return nil
}

// AddShift agrega un nuevo turno
func (s *Storage) AddShift(shift models.Shift) error {
	data, err := s.Load()
	if err != nil {
		return err
	}

	data.Shifts = append(data.Shifts, shift)

	return s.Save(data)
}

// AddEvent agrega un nuevo evento
func (s *Storage) AddEvent(event models.Event) error {
	data, err := s.Load()
	if err != nil {
		return err
	}

	data.Events = append(data.Events, event)

	return s.Save(data)
}

// GetDataPath retorna la ruta del archivo de datos
func (s *Storage) GetDataPath() string {
	return s.dataPath
}

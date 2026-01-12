package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Exercise representa un ejercicio registrado
type Exercise struct {
	Name        string    `json:"name"`
	Weight      float64   `json:"weight"`
	Reps        int       `json:"reps"`
	Sets        int       `json:"sets"`
	Date        time.Time `json:"date"`
	OneRepMax   float64   `json:"one_rep_max"`
}

// GymData contiene todos los ejercicios registrados
type GymData struct {
	Exercises []Exercise `json:"exercises"`
}

const dataFile = "gym_data.json"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "add":
		handleAdd()
	case "history":
		handleHistory()
	case "1rm":
		handle1RM()
	case "list":
		handleList()
	default:
		fmt.Printf("Comando desconocido: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Gym Tracker - Aplicación CLI para trackear ejercicios")
	fmt.Println("\nUso:")
	fmt.Println("  gym-tracker add -name <ejercicio> -weight <peso> -reps <repeticiones> -sets <series>")
	fmt.Println("  gym-tracker history -name <ejercicio>")
	fmt.Println("  gym-tracker 1rm -name <ejercicio>")
	fmt.Println("  gym-tracker list")
	fmt.Println("\nComandos:")
	fmt.Println("  add      Registrar un nuevo ejercicio")
	fmt.Println("  history  Ver el historial de un ejercicio específico")
	fmt.Println("  1rm      Calcular el 1RM (One Rep Max) de un ejercicio")
	fmt.Println("  list     Listar todos los ejercicios únicos registrados")
}

func handleAdd() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	name := addCmd.String("name", "", "Nombre del ejercicio")
	weight := addCmd.Float64("weight", 0, "Peso levantado (kg)")
	reps := addCmd.Int("reps", 0, "Número de repeticiones")
	sets := addCmd.Int("sets", 0, "Número de series")

	addCmd.Parse(os.Args[2:])

	if *name == "" || *weight <= 0 || *reps <= 0 || *sets <= 0 {
		fmt.Println("Error: Todos los parámetros son requeridos y deben ser mayores a 0")
		fmt.Println("\nUso: gym-tracker add -name <ejercicio> -weight <peso> -reps <repeticiones> -sets <series>")
		os.Exit(1)
	}

	// Calcular 1RM usando la fórmula de Epley
	oneRepMax := calculate1RM(*weight, *reps)

	exercise := Exercise{
		Name:      *name,
		Weight:    *weight,
		Reps:      *reps,
		Sets:      *sets,
		Date:      time.Now(),
		OneRepMax: oneRepMax,
	}

	data := loadData()
	data.Exercises = append(data.Exercises, exercise)
	saveData(data)

	fmt.Printf("✓ Ejercicio registrado exitosamente!\n")
	fmt.Printf("  Ejercicio: %s\n", exercise.Name)
	fmt.Printf("  Peso: %.1f kg\n", exercise.Weight)
	fmt.Printf("  Repeticiones: %d\n", exercise.Reps)
	fmt.Printf("  Series: %d\n", exercise.Sets)
	fmt.Printf("  1RM estimado: %.1f kg\n", exercise.OneRepMax)
	fmt.Printf("  Fecha: %s\n", exercise.Date.Format("2006-01-02 15:04:05"))
}

func handleHistory() {
	historyCmd := flag.NewFlagSet("history", flag.ExitOnError)
	name := historyCmd.String("name", "", "Nombre del ejercicio")

	historyCmd.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: El nombre del ejercicio es requerido")
		fmt.Println("\nUso: gym-tracker history -name <ejercicio>")
		os.Exit(1)
	}

	data := loadData()
	exercises := filterByName(data.Exercises, *name)

	if len(exercises) == 0 {
		fmt.Printf("No se encontraron registros para el ejercicio: %s\n", *name)
		return
	}

	// Ordenar por fecha (más reciente primero)
	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].Date.After(exercises[j].Date)
	})

	fmt.Printf("\n=== Historial de %s ===\n\n", *name)
	for i, ex := range exercises {
		fmt.Printf("%d. Fecha: %s\n", i+1, ex.Date.Format("2006-01-02 15:04:05"))
		fmt.Printf("   Peso: %.1f kg | Reps: %d | Sets: %d | 1RM: %.1f kg\n",
			ex.Weight, ex.Reps, ex.Sets, ex.OneRepMax)
		fmt.Println()
	}

	fmt.Printf("Total de entrenamientos: %d\n", len(exercises))
}

func handle1RM() {
	rmCmd := flag.NewFlagSet("1rm", flag.ExitOnError)
	name := rmCmd.String("name", "", "Nombre del ejercicio")

	rmCmd.Parse(os.Args[2:])

	if *name == "" {
		fmt.Println("Error: El nombre del ejercicio es requerido")
		fmt.Println("\nUso: gym-tracker 1rm -name <ejercicio>")
		os.Exit(1)
	}

	data := loadData()
	exercises := filterByName(data.Exercises, *name)

	if len(exercises) == 0 {
		fmt.Printf("No se encontraron registros para el ejercicio: %s\n", *name)
		return
	}

	// Encontrar el mejor 1RM
	maxRM := 0.0
	var bestExercise Exercise
	for _, ex := range exercises {
		if ex.OneRepMax > maxRM {
			maxRM = ex.OneRepMax
			bestExercise = ex
		}
	}

	// Calcular promedio
	totalRM := 0.0
	for _, ex := range exercises {
		totalRM += ex.OneRepMax
	}
	avgRM := totalRM / float64(len(exercises))

	fmt.Printf("\n=== 1RM para %s ===\n\n", *name)
	fmt.Printf("Mejor 1RM: %.1f kg\n", maxRM)
	fmt.Printf("  Alcanzado con: %.1f kg x %d reps\n", bestExercise.Weight, bestExercise.Reps)
	fmt.Printf("  Fecha: %s\n", bestExercise.Date.Format("2006-01-02 15:04:05"))
	fmt.Printf("\nPromedio 1RM: %.1f kg\n", avgRM)
	fmt.Printf("Total de registros: %d\n", len(exercises))
}

func handleList() {
	data := loadData()

	if len(data.Exercises) == 0 {
		fmt.Println("No hay ejercicios registrados aún")
		return
	}

	// Obtener ejercicios únicos y contar entrenamientos
	exerciseMap := make(map[string]int)
	for _, ex := range data.Exercises {
		exerciseMap[ex.Name]++
	}

	// Ordenar alfabéticamente
	names := make([]string, 0, len(exerciseMap))
	for name := range exerciseMap {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("\n=== Ejercicios Registrados ===\n")
	for i, name := range names {
		count := exerciseMap[name]
		fmt.Printf("%d. %s (%d entrenamientos)\n", i+1, name, count)
	}
	fmt.Printf("\nTotal de ejercicios únicos: %d\n", len(names))
	fmt.Printf("Total de entrenamientos: %d\n", len(data.Exercises))
}

// calculate1RM calcula el One Rep Max usando la fórmula de Epley
// 1RM = weight * (1 + reps/30)
func calculate1RM(weight float64, reps int) float64 {
	if reps == 1 {
		return weight
	}
	oneRM := weight * (1 + float64(reps)/30.0)
	return math.Round(oneRM*10) / 10 // Redondear a 1 decimal
}

func filterByName(exercises []Exercise, name string) []Exercise {
	var filtered []Exercise
	for _, ex := range exercises {
		if ex.Name == name {
			filtered = append(filtered, ex)
		}
	}
	return filtered
}

func loadData() GymData {
	data := GymData{
		Exercises: []Exercise{},
	}

	// Buscar el archivo en el directorio home o en el directorio actual
	filePath := getDataFilePath()

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		// Si el archivo no existe, devolver datos vacíos
		if os.IsNotExist(err) {
			return data
		}
		fmt.Printf("Error al leer el archivo de datos: %v\n", err)
		os.Exit(1)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Printf("Error al parsear el archivo de datos: %v\n", err)
		os.Exit(1)
	}

	return data
}

func saveData(data GymData) {
	filePath := getDataFilePath()

	// Crear directorio si no existe
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Error al crear directorio: %v\n", err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error al convertir datos a JSON: %v\n", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error al guardar el archivo de datos: %v\n", err)
		os.Exit(1)
	}
}

func getDataFilePath() string {
	// Si hay una variable de entorno DATA_DIR, usar esa ubicación
	if dataDir := os.Getenv("DATA_DIR"); dataDir != "" {
		return filepath.Join(dataDir, dataFile)
	}

	// Intentar usar el directorio home
	home, err := os.UserHomeDir()
	if err == nil {
		return filepath.Join(home, ".gym-tracker", dataFile)
	}

	// Si no se puede obtener el home, usar el directorio actual
	return dataFile
}

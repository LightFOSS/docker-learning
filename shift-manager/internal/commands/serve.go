package commands

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/LightFOSS/docker-learning/shift-manager/internal/models"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/storage"
	"github.com/LightFOSS/docker-learning/shift-manager/internal/utils"
)

const htmlTemplate = `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Shift Manager - Gestión de Turnos</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
        }

        header {
            background: white;
            padding: 30px;
            border-radius: 15px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
            margin-bottom: 30px;
            text-align: center;
        }

        h1 {
            color: #667eea;
            font-size: 2.5em;
            margin-bottom: 10px;
        }

        .subtitle {
            color: #666;
            font-size: 1.1em;
        }

        .week-navigation {
            display: flex;
            justify-content: space-between;
            align-items: center;
            background: white;
            padding: 20px;
            border-radius: 15px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
            margin-bottom: 20px;
        }

        .week-navigation h2 {
            color: #333;
            flex: 1;
            text-align: center;
        }

        .nav-button {
            background: #667eea;
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 8px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 600;
            transition: all 0.3s ease;
            text-decoration: none;
            display: inline-block;
        }

        .nav-button:hover {
            background: #5568d3;
            transform: translateY(-2px);
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.4);
        }

        .main-content {
            display: grid;
            grid-template-columns: 2fr 1fr;
            gap: 20px;
            margin-bottom: 20px;
        }

        .calendar-section, .stats-section, .list-section {
            background: white;
            padding: 25px;
            border-radius: 15px;
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
        }

        .list-section {
            grid-column: 1 / -1;
        }

        .section-title {
            color: #333;
            font-size: 1.5em;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 3px solid #667eea;
        }

        .calendar {
            display: grid;
            grid-template-columns: repeat(7, 1fr);
            gap: 10px;
        }

        .day-card {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 10px;
            min-height: 150px;
            border: 2px solid #e9ecef;
            transition: all 0.3s ease;
        }

        .day-card.today {
            border-color: #667eea;
            background: linear-gradient(135deg, #667eea15 0%, #764ba215 100%);
        }

        .day-card:hover {
            transform: translateY(-3px);
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
        }

        .day-header {
            font-weight: 700;
            color: #333;
            margin-bottom: 5px;
            font-size: 0.9em;
        }

        .day-date {
            font-size: 0.85em;
            color: #666;
            margin-bottom: 10px;
        }

        .day-card.today .day-header {
            color: #667eea;
        }

        .activity {
            background: white;
            padding: 8px;
            border-radius: 6px;
            margin-bottom: 8px;
            font-size: 0.85em;
            border-left: 3px solid #667eea;
        }

        .activity.shift {
            border-left-color: #667eea;
        }

        .activity.event {
            border-left-color: #f59e0b;
        }

        .activity-time {
            font-weight: 600;
            color: #333;
        }

        .activity-details {
            color: #666;
            margin-top: 3px;
        }

        .activity-type {
            display: inline-block;
            padding: 2px 6px;
            border-radius: 4px;
            font-size: 0.75em;
            font-weight: 600;
            margin-top: 5px;
        }

        .activity-type.normal {
            background: #dbeafe;
            color: #1e40af;
        }

        .activity-type.partido {
            background: #fce7f3;
            color: #be185d;
        }

        .activity-type.gym, .activity-type.terapia, .activity-type.clases {
            background: #fef3c7;
            color: #92400e;
        }

        .stats-grid {
            display: grid;
            gap: 15px;
        }

        .stat-card {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            border-radius: 10px;
            color: white;
        }

        .stat-label {
            font-size: 0.9em;
            opacity: 0.9;
            margin-bottom: 5px;
        }

        .stat-value {
            font-size: 2em;
            font-weight: 700;
        }

        .stat-breakdown {
            margin-top: 15px;
            padding-top: 15px;
            border-top: 1px solid rgba(255, 255, 255, 0.3);
        }

        .stat-item {
            display: flex;
            justify-content: space-between;
            margin-bottom: 8px;
            font-size: 0.9em;
        }

        .list-grid {
            display: grid;
            gap: 15px;
        }

        .list-day {
            border: 1px solid #e9ecef;
            border-radius: 10px;
            overflow: hidden;
        }

        .list-day-header {
            background: #f8f9fa;
            padding: 15px;
            font-weight: 700;
            color: #333;
            border-bottom: 2px solid #667eea;
        }

        .list-items {
            padding: 15px;
        }

        .list-item {
            display: grid;
            grid-template-columns: auto 1fr auto;
            gap: 15px;
            padding: 12px;
            background: #f8f9fa;
            border-radius: 8px;
            margin-bottom: 10px;
            align-items: center;
        }

        .list-item:last-child {
            margin-bottom: 0;
        }

        .list-item-type {
            font-weight: 700;
            padding: 5px 10px;
            border-radius: 6px;
            font-size: 0.85em;
        }

        .list-item-type.shift {
            background: #667eea;
            color: white;
        }

        .list-item-type.event {
            background: #f59e0b;
            color: white;
        }

        .list-item-info {
            display: flex;
            flex-direction: column;
            gap: 5px;
        }

        .list-item-title {
            font-weight: 600;
            color: #333;
        }

        .list-item-subtitle {
            color: #666;
            font-size: 0.9em;
        }

        .list-item-time {
            font-weight: 600;
            color: #667eea;
            white-space: nowrap;
        }

        .empty-state {
            text-align: center;
            padding: 40px;
            color: #999;
        }

        .empty-state-icon {
            font-size: 3em;
            margin-bottom: 10px;
        }

        @media (max-width: 768px) {
            .main-content {
                grid-template-columns: 1fr;
            }

            .calendar {
                grid-template-columns: 1fr;
            }

            h1 {
                font-size: 1.8em;
            }

            .list-item {
                grid-template-columns: 1fr;
                gap: 10px;
            }
        }

        .footer {
            text-align: center;
            padding: 20px;
            color: white;
            margin-top: 30px;
        }

        .footer a {
            color: white;
            text-decoration: none;
            font-weight: 600;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>📅 Shift Manager</h1>
            <p class="subtitle">Gestión de Turnos y Eventos Personales</p>
        </header>

        <div class="week-navigation">
            <a href="/?offset={{.PrevOffset}}" class="nav-button">← Semana Anterior</a>
            <h2>{{.WeekTitle}}</h2>
            <a href="/?offset={{.NextOffset}}" class="nav-button">Semana Siguiente →</a>
        </div>

        <div class="main-content">
            <div class="calendar-section">
                <h3 class="section-title">Calendario Semanal</h3>
                <div class="calendar">
                    {{range .Days}}
                    <div class="day-card{{if .IsToday}} today{{end}}">
                        <div class="day-header">{{.Weekday}}</div>
                        <div class="day-date">{{.Date}}</div>
                        {{if .Activities}}
                            {{range .Activities}}
                            <div class="activity {{.Type}}">
                                <div class="activity-time">{{.Time}}</div>
                                <div class="activity-details">{{.Details}}</div>
                                <span class="activity-type {{.SubType}}">{{.SubType}}</span>
                            </div>
                            {{end}}
                        {{else}}
                            <div style="color: #999; font-size: 0.85em; margin-top: 10px;">Sin actividades</div>
                        {{end}}
                    </div>
                    {{end}}
                </div>
            </div>

            <div class="stats-section">
                <h3 class="section-title">Estadísticas</h3>
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-label">Total Horas</div>
                        <div class="stat-value">{{.Stats.TotalHours}}</div>
                        <div class="stat-breakdown">
                            <div class="stat-item">
                                <span>Turnos</span>
                                <span>{{.Stats.ShiftCount}}</span>
                            </div>
                            <div class="stat-item">
                                <span>Eventos</span>
                                <span>{{.Stats.EventCount}}</span>
                            </div>
                        </div>
                    </div>
                    <div class="stat-card" style="background: linear-gradient(135deg, #f59e0b 0%, #d97706 100%);">
                        <div class="stat-label">Promedio Diario</div>
                        <div class="stat-value">{{.Stats.AvgDaily}}</div>
                        <div class="stat-breakdown">
                            <div class="stat-item">
                                <span>Normal</span>
                                <span>{{.Stats.NormalCount}}</span>
                            </div>
                            <div class="stat-item">
                                <span>Partido</span>
                                <span>{{.Stats.PartidoCount}}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="list-section">
            <h3 class="section-title">Lista de Actividades</h3>
            {{if .ListItems}}
                <div class="list-grid">
                    {{range .ListItems}}
                    <div class="list-day">
                        <div class="list-day-header">{{.Date}} ({{.Weekday}})</div>
                        <div class="list-items">
                            {{range .Items}}
                            <div class="list-item">
                                <div class="list-item-type {{.Type}}">{{.TypeLabel}}</div>
                                <div class="list-item-info">
                                    <div class="list-item-title">{{.Title}}</div>
                                    <div class="list-item-subtitle">{{.Subtitle}}</div>
                                </div>
                                <div class="list-item-time">{{.Time}}</div>
                            </div>
                            {{end}}
                        </div>
                    </div>
                    {{end}}
                </div>
            {{else}}
                <div class="empty-state">
                    <div class="empty-state-icon">📭</div>
                    <p>No hay actividades programadas para esta semana</p>
                </div>
            {{end}}
        </div>

        <div class="footer">
            <p>Shift Manager v1.0 | Creado con ❤️ usando Go</p>
        </div>
    </div>
</body>
</html>`

// Activity representa una actividad en el calendario
type Activity struct {
	Type    string
	SubType string
	Time    string
	Details string
}

// DayData representa los datos de un día en el calendario
type DayData struct {
	Weekday    string
	Date       string
	IsToday    bool
	Activities []Activity
}

// WebStats representa las estadísticas para la vista web
type WebStats struct {
	TotalHours    string
	AvgDaily      string
	ShiftCount    int
	EventCount    int
	NormalCount   int
	PartidoCount  int
}

// ListItem representa un item en la lista
type ListItem struct {
	Type      string
	TypeLabel string
	Title     string
	Subtitle  string
	Time      string
}

// ListDay representa un día en la lista
type ListDay struct {
	Date    string
	Weekday string
	Items   []ListItem
}

// PageData contiene todos los datos para la página
type PageData struct {
	WeekTitle  string
	PrevOffset int
	NextOffset int
	Days       []DayData
	Stats      WebStats
	ListItems  []ListDay
}

// Serve inicia el servidor web
func Serve(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.String("port", "8080", "Puerto del servidor")

	if err := fs.Parse(args); err != nil {
		return err
	}

	http.HandleFunc("/", handleIndex)

	addr := fmt.Sprintf(":%s", *port)
	fmt.Printf("🌐 Servidor web iniciado en http://localhost%s\n", addr)
	fmt.Println("   Presiona Ctrl+C para detener el servidor")

	return http.ListenAndServe(addr, nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// Obtener offset de la query string
	offsetStr := r.URL.Query().Get("offset")
	offset := 0
	if offsetStr != "" {
		if val, err := strconv.Atoi(offsetStr); err == nil {
			offset = val
		}
	}

	// Cargar datos
	store, err := storage.New()
	if err != nil {
		http.Error(w, "Error al cargar almacenamiento", http.StatusInternalServerError)
		log.Printf("Error loading storage: %v", err)
		return
	}

	data, err := store.Load()
	if err != nil {
		http.Error(w, "Error al cargar datos", http.StatusInternalServerError)
		log.Printf("Error loading data: %v", err)
		return
	}

	// Calcular la semana a mostrar
	refDate := time.Now().AddDate(0, 0, offset*7)
	startOfWeek, endOfWeek := utils.GetWeekBounds(refDate)

	// Preparar datos de la página
	pageData := PageData{
		WeekTitle:  fmt.Sprintf("Semana del %s al %s", utils.FormatDate(startOfWeek), utils.FormatDate(endOfWeek)),
		PrevOffset: offset - 1,
		NextOffset: offset + 1,
		Days:       buildDaysData(data, startOfWeek),
		Stats:      buildStats(data, startOfWeek, endOfWeek),
		ListItems:  buildListItems(data, startOfWeek, endOfWeek),
	}

	// Parsear y ejecutar template
	tmpl, err := template.New("index").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Error al cargar template", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}

	if err := tmpl.Execute(w, pageData); err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func buildDaysData(data *models.Data, startOfWeek time.Time) []DayData {
	days := make([]DayData, 7)
	today := utils.FormatDate(time.Now())

	weekdayNames := []string{"Lunes", "Martes", "Miércoles", "Jueves", "Viernes", "Sábado", "Domingo"}

	for i := 0; i < 7; i++ {
		currentDay := startOfWeek.AddDate(0, 0, i)
		dayStr := utils.FormatDate(currentDay)

		days[i] = DayData{
			Weekday:    weekdayNames[i],
			Date:       dayStr,
			IsToday:    dayStr == today,
			Activities: []Activity{},
		}

		// Agregar turnos
		for _, shift := range data.Shifts {
			if shift.Date == dayStr {
				hours, _ := shift.CalculateHours()
				activity := Activity{
					Type:    "shift",
					SubType: string(shift.Type),
					Time:    fmt.Sprintf("%s - %s", shift.StartTime, shift.EndTime),
					Details: fmt.Sprintf("Turno %s - %s", shift.Type, utils.FormatDuration(hours)),
				}
				days[i].Activities = append(days[i].Activities, activity)
			}
		}

		// Agregar eventos
		for _, event := range data.Events {
			if event.Date == dayStr {
				activity := Activity{
					Type:    "event",
					SubType: event.Type,
					Time:    fmt.Sprintf("%s - %s", event.StartTime, event.EndTime),
					Details: fmt.Sprintf("%s", event.Title),
				}
				days[i].Activities = append(days[i].Activities, activity)
			}
		}
	}

	return days
}

func buildStats(data *models.Data, startOfWeek, endOfWeek time.Time) WebStats {
	var totalHours float64
	var normalCount, partidoCount, shiftCount, eventCount int

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
			normalCount++
		} else {
			partidoCount++
		}
	}

	for _, event := range data.Events {
		eventDate, err := utils.ParseDate(event.Date)
		if err != nil {
			continue
		}

		if eventDate.Before(startOfWeek) || eventDate.After(endOfWeek) {
			continue
		}

		eventCount++
	}

	avgDaily := totalHours / 7.0

	return WebStats{
		TotalHours:   utils.FormatDuration(totalHours),
		AvgDaily:     utils.FormatDuration(avgDaily),
		ShiftCount:   shiftCount,
		EventCount:   eventCount,
		NormalCount:  normalCount,
		PartidoCount: partidoCount,
	}
}

func buildListItems(data *models.Data, startOfWeek, endOfWeek time.Time) []ListDay {
	dayMap := make(map[string]*ListDay)
	weekdayNames := map[string]string{
		"Monday":    "Lun",
		"Tuesday":   "Mar",
		"Wednesday": "Mié",
		"Thursday":  "Jue",
		"Friday":    "Vie",
		"Saturday":  "Sáb",
		"Sunday":    "Dom",
	}

	// Agregar turnos
	for _, shift := range data.Shifts {
		shiftDate, err := utils.ParseDate(shift.Date)
		if err != nil {
			continue
		}

		if shiftDate.Before(startOfWeek) || shiftDate.After(endOfWeek) {
			continue
		}

		dateStr := shift.Date
		if _, exists := dayMap[dateStr]; !exists {
			dayMap[dateStr] = &ListDay{
				Date:    dateStr,
				Weekday: weekdayNames[shiftDate.Weekday().String()],
				Items:   []ListItem{},
			}
		}

		hours, _ := shift.CalculateHours()
		item := ListItem{
			Type:      "shift",
			TypeLabel: "TURNO",
			Title:     fmt.Sprintf("Turno %s", shift.Type),
			Subtitle:  utils.FormatDuration(hours),
			Time:      fmt.Sprintf("%s - %s", shift.StartTime, shift.EndTime),
		}
		if shift.Notes != "" {
			item.Subtitle += fmt.Sprintf(" | %s", shift.Notes)
		}

		dayMap[dateStr].Items = append(dayMap[dateStr].Items, item)
	}

	// Agregar eventos
	for _, event := range data.Events {
		eventDate, err := utils.ParseDate(event.Date)
		if err != nil {
			continue
		}

		if eventDate.Before(startOfWeek) || eventDate.After(endOfWeek) {
			continue
		}

		dateStr := event.Date
		if _, exists := dayMap[dateStr]; !exists {
			dayMap[dateStr] = &ListDay{
				Date:    dateStr,
				Weekday: weekdayNames[eventDate.Weekday().String()],
				Items:   []ListItem{},
			}
		}

		item := ListItem{
			Type:      "event",
			TypeLabel: "EVENTO",
			Title:     event.Title,
			Subtitle:  event.Type,
			Time:      fmt.Sprintf("%s - %s", event.StartTime, event.EndTime),
		}
		if event.Notes != "" {
			item.Subtitle += fmt.Sprintf(" | %s", event.Notes)
		}

		dayMap[dateStr].Items = append(dayMap[dateStr].Items, item)
	}

	// Convertir map a slice ordenado
	var listDays []ListDay
	for d := 0; d < 7; d++ {
		currentDay := startOfWeek.AddDate(0, 0, d)
		dateStr := utils.FormatDate(currentDay)
		if day, exists := dayMap[dateStr]; exists {
			listDays = append(listDays, *day)
		}
	}

	return listDays
}

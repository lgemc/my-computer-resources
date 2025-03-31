package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type ResourceUsage struct {
	Date      string  `json:"date"`
	CPUUsage  float64 `json:"cpu_usage"`
	RAMUsage  float64 `json:"ram_usage"`
	DiskUsage float64 `json:"disk_usage"`
	TotalDisk float64 `json:"total_disk"`
}

var (
	db *sql.DB
	mu sync.Mutex
)

// Conectar a la base de datos
func connectDB() {
	var err error
	connStr := "postgres://danielmijares:123456@localhost:5432/system_monitor?sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a PostgreSQL:", err)
	}
	fmt.Println("✅ Conexión a PostgreSQL exitosa")
}

// Obtener el uso de CPU, RAM y Disco
func getResourceUsage() ResourceUsage {
	cpuPercent, _ := cpu.Percent(0, false)
	vmStat, _ := mem.VirtualMemory()
	diskStat, _ := disk.Usage("/")

	return ResourceUsage{
		Date:      time.Now().Format("2006-01-02 15:04:05.00000"),
		CPUUsage:  cpuPercent[0] / 100,
		RAMUsage:  float64(vmStat.Used) / float64(vmStat.Total),
		DiskUsage: float64(diskStat.Used) / float64(diskStat.Total),
		TotalDisk: float64(diskStat.Total) / (1024 * 1024 * 1024), // Convertimos a GB
	}
}

// Guardar en la base de datos
func saveToDB(usage ResourceUsage) {
	mu.Lock()
	defer mu.Unlock()

	_, err := db.Exec("INSERT INTO resource_usage (cpu_usage, ram_usage, disk_usage) VALUES ($1, $2, $3)",
		usage.CPUUsage, usage.RAMUsage, usage.DiskUsage)

	if err != nil {
		log.Println("Error insertando en PostgreSQL:", err)
	}
}

// Iniciar grabación de datos cada 30 segundos
func startRecording() {
	ticker := time.NewTicker(30 * time.Second)
	go func() {
		for range ticker.C {
			usage := getResourceUsage()
			saveToDB(usage)
		}
	}()
}

// Endpoint /resources/now
func resourceNowHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	usage := getResourceUsage()
	json.NewEncoder(w).Encode(usage)
}

// Endpoint /resources/history
func resourceHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT date, cpu_usage, ram_usage, disk_usage FROM resource_usage ORDER BY date DESC")
	if err != nil {
		http.Error(w, "Error obteniendo el historial", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Obtener el tamaño total del disco
	diskStat, _ := disk.Usage("/")
	totalDisk := float64(diskStat.Total) / (1024 * 1024 * 1024) // Convertir a GB

	var history []ResourceUsage
	for rows.Next() {
		var usage ResourceUsage
		var date time.Time

		if err := rows.Scan(&date, &usage.CPUUsage, &usage.RAMUsage, &usage.DiskUsage); err != nil {
			http.Error(w, "Error leyendo los datos", http.StatusInternalServerError)
			return
		}
		usage.TotalDisk = totalDisk                           // Asignar el valor corregido
		usage.Date = date.Format("2006-01-02 15:04:05.00000") // Formato correcto

		history = append(history, usage)
	}

	json.NewEncoder(w).Encode(history)
}

func main() {
	connectDB()
	startRecording()

	http.HandleFunc("/resources/now", resourceNowHandler)
	http.HandleFunc("/resources/history", resourceHistoryHandler)

	port := 8080
	fmt.Printf("Servidor corriendo en http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

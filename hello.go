package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Структура для JSON-ответа
type StatsResponse struct {
	CPUUsagePercent float64 `json:"cpu_usage_percent"`
	MemoryFreeMB    uint64  `json:"memory_free_mb"`
	MemoryTotalMB   uint64  `json:"memory_total_mb"`
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, "hello world from server")
}

// Новый хендлер для проверки системных ресурсов
func checkHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/check" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// 1. Получаем данные по оперативной памяти
	vMem, err := mem.VirtualMemory()
	if err != nil {
		http.Error(w, "Ошибка получения данных ОЗУ", http.StatusInternalServerError)
		return
	}

	// 2. Получаем загрузку процессора (замеряем за интервал в 200 миллисекунд)
	cpuPercentages, err := cpu.Percent(200*time.Millisecond, false)
	if err != nil || len(cpuPercentages) == 0 {
		http.Error(w, "Ошибка получения данных ЦПУ", http.StatusInternalServerError)
		return
	}

	// 3. Формируем структуру ответа (переводим байты в Мегабайты)
	response := StatsResponse{
		CPUUsagePercent: cpuPercentages[0],
		MemoryFreeMB:    vMem.Free / 1024 / 1024,
		MemoryTotalMB:   vMem.Total / 1024 / 1024,
	}

	// 4. Отдаем заголовки, что это JSON, и отправляем его
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/check", checkHandler) // Регистрируем новый эндпоинт

	log.Println("Сервер запущен на http://localhost:4321")
	log.Println("Метрики доступны на http://localhost:4321/check")

	if err := http.ListenAndServe(":4321", nil); err != nil {
		log.Fatal("Ошибка запуска:", err)
	}
}

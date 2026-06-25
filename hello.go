package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os" // Пакет для чтения переменных окружения
	"time"

	"github.com/robfig/cron/v3"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// Глобальные переменные для токенов
var (
	botToken string
	chatID   string
)

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

func checkHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/check" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	vMem, _ := mem.VirtualMemory()
	cpuPercentages, _ := cpu.Percent(200*time.Millisecond, false)

	response := StatsResponse{
		CPUUsagePercent: cpuPercentages[0],
		MemoryFreeMB:    vMem.Free / 1024 / 1024,
		MemoryTotalMB:   vMem.Total / 1024 / 1024,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func sendTelegramHeartbeat() {
	vMem, err := mem.VirtualMemory()
	if err != nil {
		log.Println("Ошибка получения ОЗУ для ТГ:", err)
		return
	}

	cpuPercentages, err := cpu.Percent(200*time.Millisecond, false)
	if err != nil || len(cpuPercentages) == 0 {
		log.Println("Ошибка получения ЦПУ для ТГ:", err)
		return
	}

	freeMB := vMem.Free / 1024 / 1024
	cpuUsage := cpuPercentages[0]

	messageText := fmt.Sprintf(
		"🟢 *Heartbeat: Сервер жив!*\n\n🖥 *ЦПУ:* %.2f%%\n💾 *Свободно ОЗУ:* %d МБ",
		cpuUsage,
		freeMB,
	)

	encodedText := url.QueryEscape(messageText)

	tgURL := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s&parse_mode=Markdown",
		botToken,
		chatID,
		encodedText,
	)

	resp, err := http.Get(tgURL)
	if err != nil {
		log.Println("Не удалось связаться с Telegram API:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Println("Сообщение успешно отправлено в Telegram!")
	} else {
		log.Printf("Telegram вернул ошибку, статус-код: %d\n", resp.StatusCode)
	}
}

func main() {
	// Считываем секреты, которые Docker передаст нам в контейнер
	botToken = os.Getenv("TG_BOT_TOKEN")
	chatID = os.Getenv("TG_CHAT_ID")

	// Проверяем, что они не пустые
	if botToken == "" || chatID == "" {
		log.Fatal("КРИТИЧЕСКАЯ ОШИБКА: Переменные TG_BOT_TOKEN или TG_CHAT_ID не заданы в окружении!")
	}

	c := cron.New()
	// Раз в 10 минут
	_, err := c.AddFunc("*/10 * * * *", func() {
		log.Println("Сработал внутренний Cron, отправляем отчет...")
		sendTelegramHeartbeat()
	})
	if err != nil {
		log.Fatal("Ошибка настройки cron:", err)
	}
	c.Start()

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/check", checkHandler)

	log.Println("Сервер запущен на порту :4321")
	log.Println("Встроенный Cron запущен на интервал раз в 10 минут.")

	if err := http.ListenAndServe(":4321", nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

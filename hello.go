package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Проверяем, что пришли именно по адресу "/hello"
	if r.URL.Path != "/hello" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}

	// 2. Отсекаем все методы, кроме GET
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// 3. Отправляем заветную строку
	fmt.Fprint(w, "hello world from server")
}

func main() {
	// Регистрируем наш хендлер
	http.HandleFunc("/hello", helloHandler)

	log.Println("Сервер запущен на http://localhost:4321/hello")

	// Запускаем сервер на порту 8080
	if err := http.ListenAndServe(":4321", nil); err != nil {
		log.Fatal("Ошибка запуска:", err)
	}
}

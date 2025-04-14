package main

import (
	"log"
	"net/http"

	"project/internal/api"
	"project/internal/storage"
)

func main() {
	// Инициализируем менеджер баз данных
	dbManager := storage.NewDBManager()

	// Настраиваем обработку статических файлов
	api.SetupStaticFiles()

	// Настраиваем маршруты API
	api.SetupRoutes(dbManager)

	port := ":8080"
	log.Printf("Сервер запущен на порту %s", port)
	log.Printf("Веб-интерфейс доступен по адресу http://localhost%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

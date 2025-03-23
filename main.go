package main

import (
	"context"
	"log"
	"myapp/internal/api"
	"myapp/internal/config"
	"myapp/internal/db"
	"myapp/internal/rss"
	"net/http"
	"sync"
)

func main() {
	// Чтение конфигурации
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Создаем подключение к базе данных
	database, err := db.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer database.Close()

	// Запуск парсинга RSS-лент в горутинах с использованием WaitGroup
	var wg sync.WaitGroup
	for _, url := range cfg.Rss {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if err := rss.FetchRssFeeds(context.Background(), url, database); err != nil {
				log.Printf("Ошибка при загрузке RSS-канала из %s: %v", url, err)
			}
		}(url)
	}

	// Ожидаем завершения всех горутин
	wg.Wait()

	// Создаем и настраиваем API
	api := api.NewAPI(database)

	// Запуск веб-сервера с обработкой ошибок
	log.Println("Запуск сервера on :8080...")
	if err := http.ListenAndServe(":8080", api.Router()); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

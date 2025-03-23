package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Success(t *testing.T) {
	// Создаем временный файл конфигурации
	configContent := `{
		"rss": ["http://example.com/feed"],
		"request_period": 60,
		"database_url": "postgres://user:password@localhost/dbname"
	}`
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Загружаем конфигурацию
	cfg, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("LoadConfig вернул ошибку: %v", err)
	}

	// Проверяем значения
	if len(cfg.Rss) != 1 || cfg.Rss[0] != "http://example.com/feed" {
		t.Errorf("Поле RSS неверно: got %v", cfg.Rss)
	}
	if cfg.RequestPeriod != 60 {
		t.Errorf("RequestPeriod field is incorrect: got %v", cfg.RequestPeriod)
	}
	if cfg.DatabaseURL != "postgres://user:password@localhost/dbname" {
		t.Errorf("Поле DatabaseURL неверно: got %v", cfg.DatabaseURL)
	}
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.json")
	if err == nil {
		t.Error("LoadConfig должен возвращать ошибку для несуществующего файла")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	// Создаем временный файл с некорректным JSON
	configContent := `{ invalid json }`
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Загружаем конфигурацию
	_, err = LoadConfig(tmpFile.Name())
	if err == nil {
		t.Error("LoadConfig должен возвращать ошибку для некорректного JSON.")
	}
}

func TestLoadConfig_ValidationError(t *testing.T) {
	// Создаем временный файл с некорректной конфигурацией
	configContent := `{
		"rss": [],
		"request_period": -1,
		"database_url": ""
	}`
	tmpFile, err := os.CreateTemp("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	// Загружаем конфигурацию
	_, err = LoadConfig(tmpFile.Name())
	if err == nil {
		t.Error("LoadConfig должен возвращать ошибку для невалидной конфигурации")
	}
}

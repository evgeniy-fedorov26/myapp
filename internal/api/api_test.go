package api

import (
	"myapp/internal/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetPosts(t *testing.T) {
	// Создаем мок базы данных
	mockDB := &db.DB{}
	api := NewAPI(mockDB)

	// Создаем тестовый запрос
	req, err := http.NewRequest("GET", "/news/10", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(api.getPosts)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Обработчик вернул неправильный код состояния: got %v want %v", status, http.StatusOK)
	}
}

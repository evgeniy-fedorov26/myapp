package api

import (
	"encoding/json"
	"fmt"
	"myapp/internal/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// API - структура для работы с API.
type API struct {
	db *db.DB
	r  *mux.Router
}

// NewAPI создает новый экземпляр API.
func NewAPI(database *db.DB) *API {
	api := &API{db: database}
	api.r = mux.NewRouter()
	api.endpoints()
	return api
}

// Router возвращает маршрутизатор для API.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация методов API в маршрутизаторе запросов.
func (api *API) endpoints() {
	api.r.HandleFunc("/news/{n}", api.getPosts).Methods(http.MethodGet, http.MethodOptions)
}

// getPosts получает `n` последних новостей.
func (api *API) getPosts(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем заголовок Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Извлекаем параметр n из URL
	vars := mux.Vars(r)
	nStr := vars["n"]

	// Преобразуем параметр n из строки в int
	n, err := strconv.Atoi(nStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Неверный параметр n: %v", err), http.StatusBadRequest)
		return
	}

	// Проверяем, что n положительное число
	if n <= 0 {
		http.Error(w, "Параметр n должен быть положительным числом", http.StatusBadRequest)
		return
	}

	// Получаем последние n новостей
	posts, err := api.db.GetLastPosts(r.Context(), n)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка синхронизации новостей из базы: %v", err), http.StatusInternalServerError)
		return
	}

	// Если новостей нет, возвращаем пустой массив
	if len(posts) == 0 {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode([]struct{}{}); err != nil {
			http.Error(w, fmt.Sprintf("Ошибка кодирования пустого списка новостей.: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Отправляем результат в формате JSON
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, fmt.Sprintf("Ошибка кодирования новостей: %v", err), http.StatusInternalServerError)
	}
}

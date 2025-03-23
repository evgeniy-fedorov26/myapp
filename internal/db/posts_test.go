package db

import (
	"context"
	"testing"
	"time"
)

func TestAddPost(t *testing.T) {
	// Создаем тестовую базу данных
	connStr := "postgres://user:password@localhost/dbname?sslmode=disable"
	db, err := NewDB(connStr)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Создаем тестовую новость
	post := &Post{
		Title:   "Test Title",
		Content: "Test Content",
		PubTime: time.Now().Unix(),
		Link:    "http://example.com",
	}

	// Добавляем новость в базу данных
	_, err = db.AddPost(context.Background(), post)
	if err != nil {
		t.Errorf("error adding post: %v", err)
	}
}

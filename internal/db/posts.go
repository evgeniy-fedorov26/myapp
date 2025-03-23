package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq" // Подключение драйвера PostgreSQL
)

// DB - структура для подключения к базе данных.
type DB struct {
	Conn *sql.DB
}

// Post - структура для представления новости.
type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	PubTime int64  `json:"pub_time"` // Время в формате Unix timestamp
	Link    string `json:"link"`
}

// NewDB создает новое подключение к базе данных.
func NewDB(connStr string) (*DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Проверяем соединение
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка пинга базы данных: %w", err)
	}

	return &DB{Conn: conn}, nil
}

// GetLastPosts возвращает последние n новостей.
func (db *DB) GetLastPosts(ctx context.Context, n int) ([]Post, error) {
	query := `SELECT title, content, pub_time, link FROM posts ORDER BY pub_time DESC LIMIT $1`
	rows, err := db.Conn.QueryContext(ctx, query, n)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		var pubTime time.Time
		if err := rows.Scan(&post.Title, &post.Content, &pubTime, &post.Link); err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		post.PubTime = pubTime.Unix()
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения строк: %w", err)
	}
	return posts, nil
}

// AddPost добавляет новую новость в базу данных.
func (db *DB) AddPost(ctx context.Context, post *Post) (int64, error) {
	query := `INSERT INTO posts (title, content, pub_time, link) VALUES ($1, $2, $3, $4) RETURNING id`
	var postID int64
	err := db.Conn.QueryRowContext(ctx, query, post.Title, post.Content, time.Unix(post.PubTime, 0), post.Link).Scan(&postID)
	if err != nil {
		return 0, fmt.Errorf("error inserting post: %w", err)
	}
	return postID, nil
}

// Close закрывает соединение с базой данных.
func (db *DB) Close() error {
	if db.Conn != nil {
		return db.Conn.Close()
	}
	return nil
}

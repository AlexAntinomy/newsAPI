package storage

import "errors"

var (
	ErrPostNotFound = errors.New("404 not found")
)

type Post struct {
	ID        int    `json:"id"`
	Author    string `json:"author_name"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt int64  `json:"created_at"`
}

type Interface interface {
	Posts() ([]Post, error) // Получение всех публикаций
	AddPost(Post) error     // Создание новой публикации
	UpdatePost(Post) error  // Обновление публикации
	DeletePost(Post) error  // Удаление публикации по ID
}

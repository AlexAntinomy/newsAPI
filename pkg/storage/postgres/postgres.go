package postgres

import (
	"database/sql"
	"news/pkg/storage"

	_ "github.com/lib/pq"
)

// Store реализует storage.Interface для PostgreSQL.
type Store struct {
	db *sql.DB
}

// New создаёт соединение с PostgreSQL.
func New(connStr string) (*Store, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id SERIAL PRIMARY KEY,
            author_name TEXT NOT NULL,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_at BIGINT NOT NULL
        );
    `)
	if err != nil {
		return nil, err
	}

	return &Store{db: db}, nil
}

// Posts возвращает все публикации из БД.
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(`
        SELECT id, author_name, title, content, created_at 
        FROM posts
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []storage.Post
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Author,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// AddPost добавляет новую публикацию.
func (s *Store) AddPost(p storage.Post) error {
	_, err := s.db.Exec(`
        INSERT INTO posts (author_name, title, content, created_at)
        VALUES ($1, $2, $3, $4)
    `, p.Author, p.Title, p.Content, p.CreatedAt)
	return err
}

// UpdatePost обновляет существующую публикацию.
func (s *Store) UpdatePost(p storage.Post) error {
	_, err := s.db.Exec(`
        UPDATE posts 
        SET author_name = $1, title = $2, content = $3, created_at = $4
        WHERE id = $5
    `, p.Author, p.Title, p.Content, p.CreatedAt, p.ID)
	return err
}

// DeletePost удаляет публикацию по ID.
func (s *Store) DeletePost(p storage.Post) error {
	result, err := s.db.Exec("DELETE FROM posts WHERE id = $1", p.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return storage.ErrPostNotFound
	}

	return nil
}

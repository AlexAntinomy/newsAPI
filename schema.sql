DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_name TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT NOT NULL
);

INSERT INTO posts (id, author_name, title, content, created_at) 
VALUES 
(1, 'Дмитрий', 'Effective Go', 'Go is a new language...', 0),
(2, 'Дмитрий', 'The Go Memory Model', 'The Go memory model specifies...', 0);
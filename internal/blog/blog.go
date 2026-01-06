package blog

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Blog struct {
	ID          int
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Post struct {
	Blog
	Markdown string // markdown
}

type MarkdownDocument struct {
	Title       string
	Description string
	Filepath    string
	LastUpdated time.Time
}

type BlogTable struct {
	ID          int
	Title       string
	Description string
	Markdown    string
	CreatedAt   time.Time
}

var DB *sql.DB

func InitDatabase() error {
	var err error
	DB, err = sql.Open("sqlite", "./blog.db")
	if err != nil {
		return err
	}

	// Create table if it doesn't exist
	query := `
	CREATE TABLE IF NOT EXISTS blogs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		description TEXT,
		markdown TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = DB.Exec(query)

	return nil
}

func GetListOfBlogInfo() ([]Blog, error) {
	rows, err := DB.Query("SELECT id, title, description, created_at FROM blogs ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blogs []Blog
	for rows.Next() {
		var b Blog
		if err := rows.Scan(&b.ID, &b.Title, &b.Description, &b.CreatedAt); err != nil {
			return nil, err
		}
		blogs = append(blogs, b)
	}
	return blogs, nil
}

func GetBlogByID(id int) (*Post, error) {
	row := DB.QueryRow("SELECT id, title, description, markdown, created_at FROM blogs WHERE id = ?", id)

	var p Post
	err := row.Scan(&p.ID, &p.Title, &p.Description, &p.Markdown, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

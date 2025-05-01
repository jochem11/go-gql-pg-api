package todo

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Repository interface {
	Close()
	PutTodo(ctx context.Context, t Todo) error
	GetTodoByID(ctx context.Context, id string) (*Todo, error)
	ListTodos(ctx context.Context) ([]Todo, error)
}

type postgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (Repository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return &postgresRepository{db}, nil
}

func (r postgresRepository) Close() {
	r.Close()
}

func (r postgresRepository) PutTodo(ctx context.Context, t Todo) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO todos(id, text, completed, updated_at) VALUES ($1, $2, $3, $4)", t.ID, t.Text, t.Completed, t.UpdatedAt)
	return err
}

func (r postgresRepository) GetTodoByID(ctx context.Context, id string) (*Todo, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id, text, completed, updated_at FROM todos WHERE id = $id", id)
	t := &Todo{}
	if err := row.Scan(&t.ID, &t.Text, &t.Completed, &t.UpdatedAt); err != nil {
		return nil, err
	}
	return t, nil
}

func (r postgresRepository) ListTodos(ctx context.Context) ([]Todo, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, text,completed, updated_at FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := []Todo{}

	for rows.Next() {
		t := &Todo{}
		if err := rows.Scan(&t.ID, &t.Text, &t.Completed, &t.UpdatedAt); err == nil {
			todos = append(todos, *t)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

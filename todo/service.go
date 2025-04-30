package todo

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type Service interface {
	PostTodo(ctx context.Context, text string, completed *bool) (*Todo, error)
	GetTodo(ctx context.Context, id string) (*Todo, error)
	GetTodos(ctx context.Context) ([]Todo, error)
}

type Todo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	UpdatedAt time.Time `json:"updated_at"`
}

type todoService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &todoService{r}
}

func (s todoService) PostTodo(ctx context.Context, text string, completed *bool) (*Todo, error) {
	id := ksuid.New()
	t := &Todo{
		ID:        id.String(),
		Text:      text,
		UpdatedAt: id.Time(),
	}

	if completed != nil {
		t.Completed = *completed
	} else {
		t.Completed = false
	}

	if err := s.repository.PutTodo(ctx, *t); err != nil {
		return nil, err
	}

	return t, nil
}

func (s todoService) GetTodo(ctx context.Context, id string) (*Todo, error) {
	return s.repository.GetTodoByID(ctx, id)
}

func (s todoService) GetTodos(ctx context.Context) ([]Todo, error) {
	return s.repository.ListTodos(ctx)
}

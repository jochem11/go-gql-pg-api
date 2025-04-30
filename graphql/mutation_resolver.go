package main

import (
	"context"
	model "github.com/jochem11/go-gql-pg-api/graphql/generated/models"
	"log"
	"time"
)

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver) CreateTodo(ctx context.Context, todo model.CreateTodoInput) (*model.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	t, err := r.server.todoClient.PostTodo(ctx, todo.Text, todo.Completed)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &model.Todo{
		ID:        t.ID,
		Text:      t.Text,
		Completed: t.Completed,
		UpdatedAt: t.UpdatedAt,
	}, nil
}

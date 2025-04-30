package main

import (
	"context"
	model "github.com/jochem11/go-gql-pg-api/graphql/generated/models"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Todos(ctx context.Context, id *string) ([]*model.Todo, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.todoClient.GetTodo(ctx, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return []*model.Todo{{
			ID:        r.ID,
			Text:      r.Text,
			Completed: r.Completed,
			UpdatedAt: r.UpdatedAt,
		}}, nil
	}

	todoList, err := r.server.todoClient.GetTodos(ctx)
	var todos []*model.Todo

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, t := range todoList {
		todo := &model.Todo{
			ID:        t.ID,
			Text:      t.Text,
			Completed: t.Completed,
			UpdatedAt: t.UpdatedAt,
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

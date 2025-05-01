package todo

import (
	"context"
	_ "fmt"
	"github.com/jochem11/go-gql-pg-api/todo/pb"
	"google.golang.org/grpc"
	_ "time"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.TodoServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c := pb.NewTodoServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) PostTodo(ctx context.Context, text string, completed bool) (*Todo, error) {
	r, err := c.service.PostTodo(ctx, &pb.PostTodoRequest{
		Text:      text,
		Completed: completed,
	})
	if err != nil {
		return nil, err
	}

	// Convert google.protobuf.Timestamp to time.Time
	updatedAt := r.Todo.UpdatedAt.AsTime()

	return &Todo{
		ID:        r.Todo.Id,
		Text:      r.Todo.Text,
		Completed: r.Todo.Completed,
		UpdatedAt: updatedAt, // Use the converted time.Time
	}, nil
}

func (c *Client) GetTodo(ctx context.Context, id string) (*Todo, error) {
	r, err := c.service.GetTodo(ctx, &pb.GetTodoRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	// Convert google.protobuf.Timestamp to time.Time
	updatedAt := r.Todo.UpdatedAt.AsTime()

	return &Todo{
		ID:        r.Todo.Id,
		Text:      r.Todo.Text,
		Completed: r.Todo.Completed,
		UpdatedAt: updatedAt, // Use the converted time.Time
	}, nil
}

func (c *Client) GetTodos(ctx context.Context) ([]Todo, error) {
	r, err := c.service.GetTodos(ctx, &pb.GetTodosRequest{})
	if err != nil {
		return nil, err
	}

	todos := []Todo{}
	for _, t := range r.Todos {
		// Convert google.protobuf.Timestamp to time.Time
		updatedAt := t.UpdatedAt.AsTime()

		todos = append(todos, Todo{
			ID:        t.Id,
			Text:      t.Text,
			Completed: t.Completed,
			UpdatedAt: updatedAt, // Use the converted time.Time
		})
	}
	return todos, nil
}

package todo

import (
	"context"
	"fmt"
	"github.com/jochem11/go-gql-pg-api/todo/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type grpcServer struct {
	pb.UnimplementedTodoServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	pb.RegisterTodoServiceServer(serv, &grpcServer{
		service: s,
	})
	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostTodo(ctx context.Context, r *pb.PostTodoRequest) (*pb.PostTodoResponse, error) {
	t, err := s.service.PostTodo(ctx, r.Text, r.Completed)
	if err != nil {
		return nil, err
	}
	return &pb.PostTodoResponse{Todo: &pb.Todo{
		Id:        t.ID,
		Text:      t.Text,
		Completed: t.Completed,
	}}, nil
}

func (s *grpcServer) GetTodo(ctx context.Context, r *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	t, err := s.service.GetTodo(ctx, r.Id)
	if err != nil {
		return nil, err
	}

	updatedAtBytes, err := t.UpdatedAt.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal UpdatedAt: %v", err)
	}

	return &pb.GetTodoResponse{Todo: &pb.Todo{
		Id:        t.ID,
		Text:      t.Text,
		Completed: t.Completed,
		UpdatedAt: updatedAtBytes,
	}}, nil
}

func (s *grpcServer) GetTodos(ctx context.Context, r *pb.GetTodosRequest) (*pb.GetTodosResponse, error) {
	res, err := s.service.GetTodos(ctx)
	if err != nil {
		return nil, err
	}

	todos := []*pb.Todo{}
	for _, t := range res {
		updatedAtBytes, err := t.UpdatedAt.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal UpdatedAt for todo ID %v: %v", t.ID, err)
		}

		todos = append(
			todos,
			&pb.Todo{
				Id:        t.ID,
				Text:      t.Text,
				Completed: t.Completed,
				UpdatedAt: updatedAtBytes,
			},
		)
	}
	return &pb.GetTodosResponse{Todos: todos}, nil
}

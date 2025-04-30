package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/jochem11/go-gql-pg-api/graphql/generated"
	"github.com/jochem11/go-gql-pg-api/todo"
)

type Server struct {
	todoClient *todo.Client
}

func (s *Server) Mutation() generated.MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() generated.QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func NewGraphQLServer(todoURL string) (*Server, error) {
	todoClient, err := todo.NewClient(todoURL)
	if err != nil {
		return nil, err
	}

	return &Server{
		todoClient,
	}, nil
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return generated.NewExecutableSchema(generated.Config{
		Resolvers: s,
	})
}

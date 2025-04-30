package main

import (
	"github.com/jochem11/go-gql-pg-api/todo"
	"github.com/kelseyhightower/envconfig"
	"github.com/tinrab/retry"
	"log"
	"time"
)

type Config struct {
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	var r todo.Repository
	retry.ForeverSleep(2*time.Second, func(_ int) (err error) {
		r, err = todo.NewPostgresRepository(cfg.DatabaseURL)
		if err != nil {
			log.Println(err)
		}
		return
	})
	defer r.Close()

	log.Println("Listening on port 8080...")
	s := todo.NewService(r)
	log.Fatal(todo.ListenGRPC(s, 8080))
}

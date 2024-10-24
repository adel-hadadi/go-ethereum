package main

import (
	"context"
	"ethereum/internal/common/logs"
	"ethereum/internal/common/server"
	"ethereum/internal/ethereum/ports"
	"ethereum/internal/ethereum/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	logs.Init()

	ctx := context.Background()

	application := service.NewApplication(ctx)

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(
			ports.NewHTTPServer(application),
			router,
		)
	})
}

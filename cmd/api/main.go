package main

import (
	"ethereum/cmd/application"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

const (
	webPort = "8081"

	ErrGetAccountBalance = "Error on get account balance"
)

type BalanceRes struct {
	Balance *big.Float `json:"balance"`
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error on loading .env file")
	}

	mux := chi.NewRouter()

	mux.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	})

	mux.Use(middleware.Logger, middleware.StripSlashes)

	client, err := ethclient.Dial(os.Getenv("NETWORK_ADDRESS"))
	if err != nil {
		log.Fatal(fmt.Errorf("Error on connecting to sepolia: %w", err))
	}

	application.New(mux, client).Serve()
}

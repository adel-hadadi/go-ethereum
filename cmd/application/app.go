package application

import (
	transactionhandler "ethereum/internal/handler/transaction"
	wallethandler "ethereum/internal/handler/wallet"
	ethereumservice "ethereum/service/ethereum"
	"fmt"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-chi/chi/v5"
)

const webPort = "8081"

type Application struct {
	Mux       *chi.Mux
	EthClient *ethclient.Client
	Handlers  Handlers
	Services  Services
}

type Handlers struct {
	TransactionHandler transactionhandler.Handler
	WalletHandler      wallethandler.Handler
}

type Services struct {
	EthereumService ethereumservice.EthereumService
}

func New(mux *chi.Mux, client *ethclient.Client) *Application {
	app := &Application{
		Mux:       mux,
		EthClient: client,
	}

	app.InitServices().InitHandlers()

	return app
}

func (app *Application) InitServices() *Application {
	app.Services = Services{
		EthereumService: ethereumservice.New(app.EthClient),
	}

	return app
}

func (app *Application) InitHandlers() *Application {
	app.Handlers = Handlers{
		TransactionHandler: transactionhandler.New(app.Services.EthereumService),
		WalletHandler:      wallethandler.New(app.Services.EthereumService),
	}

	return app
}

func (app *Application) initRoutes() {
	app.Mux.Route("/api/v1", func(r chi.Router) {
		app.Handlers.TransactionHandler.Routes(r)
		app.Handlers.WalletHandler.Routes(r)
	})
}

func (app *Application) Serve() {
	app.initRoutes()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

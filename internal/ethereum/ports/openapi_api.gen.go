package ports

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ServerInterface interface {
	// Balance (GET /wallet/balance)
	Balance(w http.ResponseWriter, r *http.Request)

	// Deposit (GET /wallet/deposit)
	Deposit(w http.ResponseWriter, r *http.Request)

	// Withdraw (GET /wallet/withdraw)
	Withdraw(w http.ResponseWriter, r *http.Request)

	// TransactionCheck (POST /transactions/check)
	TransactionCheck(w http.ResponseWriter, r *http.Request)

	// CreateTransaction (POST /transaction/create)
	CreateTransaction(w http.ResponseWriter, r *http.Request)
}

type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
}

type MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc

func (siw *ServerInterfaceWrapper) Balance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Balance(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

func (siw *ServerInterfaceWrapper) TransactionCheck(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.TransactionCheck(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

func (siw *ServerInterfaceWrapper) Deposit(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Deposit(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

func (siw *ServerInterfaceWrapper) Withdraw(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Withdraw(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

func (siw *ServerInterfaceWrapper) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTransaction(w, r)
	}

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler(w, r.WithContext(ctx))
}

type ChiServerOptions struct {
	BaseURL     string
	BaseRouter  chi.Router
	Middlewares []MiddlewareFunc
}

// HandlerFromMux creates http.Handler with routing matching OpenApi spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/wallet/balance", wrapper.Balance)
	})

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/wallet/deposit", wrapper.Deposit)
	})

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/wallet/withdraw", wrapper.Withdraw)
	})

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/transactions/check", wrapper.TransactionCheck)
	})

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/transactions", wrapper.CreateTransaction)
	})

	return r
}

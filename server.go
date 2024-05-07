package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"account-test/config"
	"account-test/internal/core/services"
	"account-test/internal/repositories"
	db "account-test/postgres"
	"account-test/static"
)

func main() {
	config.InitReader()
	port := os.Getenv("PORT")
	if port == "" {
		panic(static.EmptyPort)
	}

	r := chi.NewRouter()

	// Start of Dependency Injection
	appConfig := config.Init()
	dbClient, err := db.Init(appConfig.DB)
	if err != nil {
		panic(err)
	}
	accountPort := repositories.NewAccountPort(dbClient, appConfig.DB)
	transactionPort := repositories.NewTransactionPort(dbClient, appConfig.DB)

	accountSvc := services.NewAccountSvc(accountPort)
	transactionSvc := services.NewTransactionSvc(accountPort, transactionPort)
	// End of Dependency Injection

	r.Group(func(r chi.Router) {
		r.Route("/accounts", func(route chi.Router) {
			route.Get("/{account_id}", accountSvc.GetAccount)
			route.Post("/", accountSvc.PostAccount)
		})
		r.Route("/transactions", func(route chi.Router) {
			route.Post("/", transactionSvc.PostTransaction)
		})
	})

	// Health Endpoint for Liveness Probe
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("app running on http://localhost:%s/", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

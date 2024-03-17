package main

import (
	"log"
	"net/http"

	"github.com/arosace/WellnessWaveApi/internal/account/app/handler"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
)

// ServiceSetup holds all the services and their handlers for the application.
type ServiceSetup struct {
	AccountServiceHandler *handler.AccountHandler
	// Add other service handlers here as you expand
}

func main() {
	services := initializeServices()

	registerEndpoints(services)

	// Start the HTTP server
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initializeServices sets up all the services, repositories, and handlers.
func initializeServices() ServiceSetup {
	encryptor := &encryption.Encryptor{
		Passphrase: "randompassphraseof32bytes1234567",
	}

	// Initialize all repositories
	accountRepo := repository.NewMockAccountRepository()

	// Initialize all services with their respective repositories
	accountService := service.NewAccountService(accountRepo, encryptor)

	// Initialize all handlers with their respective services
	accountServiceHandler := handler.NewAccountHandler(accountService)

	// Return a struct containing all handlers for easy access
	return ServiceSetup{
		AccountServiceHandler: accountServiceHandler,
	}
}

func registerEndpoints(services ServiceSetup) {
	http.HandleFunc("/accounts/register", services.AccountServiceHandler.HandleAddAccount)
	http.HandleFunc("/accounts", services.AccountServiceHandler.HandleGetAccounts)
	http.HandleFunc("/accounts/attach", services.AccountServiceHandler.HandleAttachAccount)
}

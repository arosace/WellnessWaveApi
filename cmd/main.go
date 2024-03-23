package main

import (
	"log"
	"net/http"

	"github.com/arosace/WellnessWaveApi/cmd/account"
	"github.com/arosace/WellnessWaveApi/cmd/event"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/gorilla/mux"
)

type Service interface {
	Init()
	RegisterEndpoints()
}

// ServiceSetup holds all the services and their handlers for the application.
type ServiceSetup struct {
	AccountService *account.AccountService
	EventService   *event.EventService
}

func main() {
	router, _ := initializeServices()

	// Start the HTTP server
	log.Println("Starting server on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initializeServices sets up all the services, repositories, and handlers.
func initializeServices() (*mux.Router, ServiceSetup) {
	//initialize router
	r := mux.NewRouter()
	//initialize encryptor
	encryptor := &encryption.Encryptor{
		Passphrase: "randompassphraseof32bytes1234567",
	}

	//initialize account service
	accServ := account.AccountService{
		Router:    r,
		Encryptor: encryptor,
	}
	accServ.Init()

	log.Println("Account service is up")

	//initialize event service
	eventServ := event.EventService{Router: r}
	eventServ.Init()

	log.Println("Event service is up")

	return r, ServiceSetup{
		AccountService: &accServ,
		EventService:   &eventServ,
	}
}

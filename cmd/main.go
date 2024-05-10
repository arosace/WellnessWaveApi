package main

import (
	"log"

	"github.com/arosace/WellnessWaveApi/cmd/account"
	"github.com/arosace/WellnessWaveApi/cmd/event"
	"github.com/arosace/WellnessWaveApi/cmd/planner"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/gorilla/mux"

	"github.com/pocketbase/pocketbase"
)

type Service interface {
	Init()
	RegisterEndpoints()
}

// ServiceSetup holds all the services and their handlers for the application.
type ServiceSetup struct {
	AccountService *account.AccountService
	EventService   *event.EventService
	PlannerService *planner.PlannerService
}

func main() {
	app := pocketbase.New()

	initializeServices(app)

	// Start the HTTP server
	log.Println("Starting server on port 8090...")
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initializeServices sets up all the services, repositories, and handlers.
func initializeServices(app *pocketbase.PocketBase) (*mux.Router, ServiceSetup) {
	//initialize router
	r := mux.NewRouter()
	//initialize encryptor
	encryptor := &encryption.Encryptor{
		Passphrase: "randompassphraseof32bytes1234567",
	}

	//initialize account service
	accServ := account.AccountService{
		App:       app,
		Router:    r,
		Encryptor: encryptor,
	}
	accServ.Init()

	log.Println("Account service is up")

	//initialize event service
	eventServ := event.EventService{Router: r}
	eventServ.Init()

	log.Println("Event service is up")

	//initialize planner service
	plannerServ := planner.PlannerService{Router: r}
	plannerServ.Init()

	log.Println("Planner service is up")

	return r, ServiceSetup{
		AccountService: &accServ,
		EventService:   &eventServ,
		PlannerService: &plannerServ,
	}
}

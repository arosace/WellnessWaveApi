package main

import (
	"log"

	"github.com/arosace/WellnessWaveApi/cmd/account"
	"github.com/arosace/WellnessWaveApi/cmd/event"
	"github.com/arosace/WellnessWaveApi/cmd/planner"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type Service interface {
	Init()
	RegisterEndpoints()
	RegisterHooks()
}

// ServiceSetup holds all the services and their handlers for the application.
type ServiceSetup struct {
	AccountService *account.AccountService
	EventService   *event.EventService
	PlannerService *planner.PlannerService
}

func main() {
	app := pocketbase.New()

	app.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		initializeServices(app)
		return nil
	})

	// Start the HTTP server
	log.Println("Starting server on port 8090...")
	if err := app.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initializeServices sets up all the services, repositories, and handlers.
func initializeServices(app *pocketbase.PocketBase) {
	// initialize mailer
	mailer := app.NewMailClient()
	// initialize dao
	dao := app.Dao()

	if mailer == nil {
		log.Fatal("mailer was not initiated properly")
	}
	if dao == nil {
		log.Fatal("dao was not initiated properly")
	}

	//initialize encryptor
	encryptor := &encryption.Encryptor{
		Passphrase: "randompassphraseof32bytes1234567",
	}

	//initialize account service
	accServ := account.AccountService{
		App:       app,
		Mailer:    mailer,
		Dao:       dao,
		Encryptor: encryptor,
	}
	accServ.Init()

	log.Println("Account service is up")

	//initialize event service
	eventServ := event.EventService{
		App: app,
		Dao: dao,
	}
	eventServ.Init()

	log.Println("Event service is up")

	//initialize planner service
	plannerServ := planner.PlannerService{App: app}
	plannerServ.Init()

	log.Println("Planner service is up")
}

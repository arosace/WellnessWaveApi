package account

import (
	"github.com/arosace/WellnessWaveApi/internal/account/handler"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type AccountService struct {
	App            *pocketbase.PocketBase
	Router         *mux.Router
	Encryptor      *encryption.Encryptor
	ServiceHandler *handler.AccountHandler
}

func (s AccountService) Init() {
	// Initialize all repositories
	accountRepo := repository.NewMockAccountRepository()

	// Initialize all services with their respective repositories
	accountService := service.NewAccountService(accountRepo, s.Encryptor)

	// Initialize all handlers with their respective services
	accountServiceHandler := handler.NewAccountHandler(accountService)

	s.ServiceHandler = accountServiceHandler

	s.RegisterEndpoints()
}

func (s AccountService) RegisterEndpoints() {
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/accounts", s.ServiceHandler.HandleGetAccounts, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/api/accounts/register", s.ServiceHandler.HandleAddAccount, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/accounts/:id", s.ServiceHandler.HandleGetAccountsById, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/api/accounts/attach", s.ServiceHandler.HandleAttachAccount, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/accounts/attached/:parent_id", s.ServiceHandler.HandleGetAttachedAccounts, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.PUT("/api/accounts/update", s.ServiceHandler.HandleUpdateAccount, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/api/accounts/login", s.ServiceHandler.HandleLogIn, utils.EchoMiddleware)
		return nil
	})
	//s.Router.HandleFunc("/api/register", utils.HttpMiddleware(s.ServiceHandler.HandleAddAccount))
	//s.Router.HandleFunc("/accounts", utils.HttpMiddleware(s.ServiceHandler.HandleGetAccounts))
	//s.Router.HandleFunc("/accounts/{id}", utils.HttpMiddleware(s.ServiceHandler.HandleGetAccountsById)).Methods("GET", "OPTIONS")
	//s.Router.HandleFunc("/accounts/attach", utils.HttpMiddleware(s.ServiceHandler.HandleAttachAccount)).Methods("POST", "OPTIONS")
	//s.Router.HandleFunc("/accounts/attached/{parent_id}", utils.HttpMiddleware(s.ServiceHandler.HandleGetAttachedAccounts))
	//s.Router.HandleFunc("/accounts/update", utils.HttpMiddleware(s.ServiceHandler.HandleUpdateAccount))
	//s.Router.HandleFunc("/login", utils.HttpMiddleware(s.ServiceHandler.HandleLogIn))
}

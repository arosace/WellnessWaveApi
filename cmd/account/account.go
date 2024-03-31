package account

import (
	"github.com/arosace/WellnessWaveApi/internal/account/handler"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/gorilla/mux"
)

type AccountService struct {
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
	s.Router.HandleFunc("/accounts/register", s.ServiceHandler.HandleAddAccount)
	s.Router.HandleFunc("/accounts", s.ServiceHandler.HandleGetAccounts)
	s.Router.HandleFunc("/accounts/attach", s.ServiceHandler.HandleAttachAccount)
	s.Router.HandleFunc("/accounts/{parent_id}", s.ServiceHandler.HandleGetAttachedAccounts)
	s.Router.HandleFunc("/accounts/update", s.ServiceHandler.HandleUpdateAccount)
}

package account

import (
	"fmt"
	"net/http"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
	"github.com/arosace/WellnessWaveApi/internal/account/handler"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type AccountService struct {
	App            *pocketbase.PocketBase
	Encryptor      *encryption.Encryptor
	ServiceHandler *handler.AccountHandler
	Mailer         mailer.Mailer
	Dao            *daos.Dao
}

func (s AccountService) Init() {
	// Initialize all repositories
	accountRepo := repository.NewAccountRepository(s.Dao)
	// Initialize all services with their respective repositories
	accountService := service.NewAccountService(accountRepo, s.Encryptor)
	// Initialize all handlers with their respective services
	accountServiceHandler := handler.NewAccountHandler(accountService)
	s.ServiceHandler = accountServiceHandler
	s.RegisterEndpoints()
	s.RegisterHooks()
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
		e.Router.PUT("/api/accounts/verify", s.ServiceHandler.HandleVerifyAccount, utils.EchoMiddleware)
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
}

func (s AccountService) RegisterHooks() {
	s.App.OnModelAfterCreate(domain.TableName).Add(func(e *core.ModelEvent) error {
		record := e.Model.(*models.Record)
		switch record.GetString("role") {
		case domain.HealthSpecialistRole:
			if err := utils.SendVerifyAccountHealthSpecialistEmail(
				s.Mailer,
				"hello@noreply.com",
				record.GetString("username"),
				record.GetString("email"),
			); err != nil {
				return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("Failed to send email:%s", err.Error()), err)
			}
		case domain.PatientRole:
			if err := utils.SendVerifyAccountPatientEmail(
				s.Mailer,
				"hello@noreply.com",
				record.GetString("username"),
				record.GetString("email"),
				record.GetString("encrypted_password"),
			); err != nil {
				return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("Failed to send email:%s", err.Error()), err)
			}
		default:
			return nil
		}
		return nil
	})
}

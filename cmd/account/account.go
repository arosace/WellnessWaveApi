package account

import (
	"fmt"
	"net/http"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
	"github.com/arosace/WellnessWaveApi/internal/account/handler"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
	eventDomain "github.com/arosace/WellnessWaveApi/internal/event/domain"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type AccountService struct {
	App                  *pocketbase.PocketBase
	Encryptor            *encryption.Encryptor
	ServiceHandler       *handler.AccountHandler
	RepositoryInteractor *repository.AccountRepo
	Mailer               mailer.Mailer
	Dao                  *daos.Dao
}

func (s AccountService) Init() {
	// Initialize all repositories
	accountRepo := repository.NewAccountRepository(s.Dao)
	// Initialize all services with their respective repositories
	accountService := service.NewAccountService(accountRepo, s.Encryptor)
	// Initialize all handlers with their respective services
	accountServiceHandler := handler.NewAccountHandler(accountService)
	s.ServiceHandler = accountServiceHandler
	s.RepositoryInteractor = accountRepo
	s.RegisterEndpoints()
	s.RegisterHooks()
}

func (s AccountService) RegisterEndpoints() {
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/v1/accounts", s.ServiceHandler.HandleGetAccounts, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/v1/accounts/register", s.ServiceHandler.HandleAddAccount, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.PUT("/v1/accounts/verify", s.ServiceHandler.HandleVerifyAccount, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/v1/accounts/:id", s.ServiceHandler.HandleGetAccountsById, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/v1/accounts/attach", s.ServiceHandler.HandleAttachAccount, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/v1/accounts/attached/:parent_id", s.ServiceHandler.HandleGetAttachedAccounts, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.PUT("/v1/accounts/update", s.ServiceHandler.HandleUpdateAccount, utils.EchoMiddleware)
		return nil
	})
	s.App.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/v1/accounts/login", s.ServiceHandler.HandleLogIn, utils.EchoMiddleware)
		return nil
	})
}

func (s AccountService) RegisterHooks() {
	// listens for changes to the "accounts" table and acts accordingly (sends an email to the newly created account)
	s.App.OnModelAfterCreate(domain.TableName).Add(func(e *core.ModelEvent) error {
		record := e.Model.(*models.Record)
		switch record.GetString("role") {
		case domain.HealthSpecialistRole:
			if err := utils.SendVerifyAccountHealthSpecialistEmail(
				s.Mailer,
				record.GetString("username"),
				record.GetString("email"),
			); err != nil {
				return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("Failed to send email:%s", err.Error()), err)
			}
		case domain.PatientRole:
			if err := utils.SendVerifyAccountPatientEmail(
				s.Mailer,
				record.GetString("username"),
				record.GetString("email"),
				record.GetString("encrypted_password"),
				record.Id,
			); err != nil {
				return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("Failed to send email:%s", err.Error()), err)
			}
		default:
			return nil
		}
		return nil
	})

	// listens for changes to the "events" table and acts accordingly (checks if health specialist id and patient id exist)
	s.App.OnModelBeforeCreate(eventDomain.TABLENAME).Add(func(e *core.ModelEvent) error {
		record := e.Model.(*models.Record)
		ctx := &echo.DefaultContext{}
		_, err := s.RepositoryInteractor.FindByID(ctx, record.GetString("health_specialist_id"))
		if err != nil {
			if utils.IsErrorNotFound(err) {
				return apis.NewApiError(http.StatusNotFound, "event was not created because health specialist id does not exist", err)
			}
			return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("there was an error verifyin that health specialist id exists:%s", err.Error()), err)
		}
		_, err = s.RepositoryInteractor.FindByID(ctx, record.GetString("patient_id"))
		if err != nil {
			if utils.IsErrorNotFound(err) {
				return apis.NewApiError(http.StatusNotFound, "event was not created because patient does not exist", err)
			}
			return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("there was an error verifyin that patient id exists:%s", err.Error()), err)
		}
		return nil
	})

	// listens for changes to the "events" table and acts accordingly (sends reminder email to patient about call)
	s.App.OnModelAfterCreate(eventDomain.TABLENAME).Add(func(e *core.ModelEvent) error {
		event := e.Model.(*models.Record)
		ctx := &echo.DefaultContext{}
		patient, err := s.RepositoryInteractor.FindByID(ctx, event.GetString("patient_id"))
		if err != nil {
			if utils.IsErrorNotFound(err) {
				return apis.NewApiError(http.StatusNotFound, "event was not created because patient does not exist", err)
			}
			return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("there was an error verifyin that patient id exists:%s", err.Error()), err)
		}

		err = utils.SendEventEmailToPatient(
			s.Mailer,
			patient.GetString("username"),
			patient.Email(),
			event,
		)
		if err != nil {
			return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("Failed to send email:%s", err.Error()), err)
		}
		return nil
	})

	// listens for updates to the "events" table and acts accordingly (sends reminder email to patient about call)
	s.App.OnModelAfterUpdate(eventDomain.TABLENAME).Add(func(e *core.ModelEvent) error {
		event := e.Model.(*models.Record)
		ctx := &echo.DefaultContext{}
		patient, err := s.RepositoryInteractor.FindByID(ctx, event.GetString("patient_id"))
		if err != nil {
			if utils.IsErrorNotFound(err) {
				return apis.NewApiError(http.StatusNotFound, "event was not created because patient does not exist", err)
			}
			return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("there was an error verifyin that patient id exists:%s", err.Error()), err)
		}

		err = utils.SendRescheduleEventEmailToPatient(
			s.Mailer,
			patient.GetString("username"),
			patient.Email(),
			event,
		)
		if err != nil {
			return apis.NewApiError(http.StatusBadRequest, fmt.Sprintf("Failed to send email:%s", err.Error()), err)
		}
		return nil
	})
}

package account

import (
	"fmt"
	"net/http"
	"net/mail"

	"github.com/arosace/WellnessWaveApi/internal/account/domain"
	"github.com/arosace/WellnessWaveApi/internal/account/handler"
	"github.com/arosace/WellnessWaveApi/internal/account/repository"
	"github.com/arosace/WellnessWaveApi/internal/account/service"
	"github.com/arosace/WellnessWaveApi/pkg/utils"
	encryption "github.com/arosace/WellnessWaveApi/pkg/utils"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/mailer"
)

type AccountService struct {
	App            *pocketbase.PocketBase
	Encryptor      *encryption.Encryptor
	ServiceHandler *handler.AccountHandler
}

func (s AccountService) Init() {
	// Initialize all repositories
	accountRepo := repository.NewAccountRepository(s.App)
	// Initialize all services with their respective repositories
	accountService := service.NewAccountService(accountRepo, s.Encryptor)
	// Initialize all handlers with their respective services
	accountServiceHandler := handler.NewAccountHandler(accountService, s.App)
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
	s.App.OnModelAfterCreate("accounts").Add(func(e *core.ModelEvent) error {
		record := e.Model.(*models.Record)
		if record.GetString("role") == domain.HhealthSpecialistRole {
			// Generate verification token
			token, err := utils.GenerateVerificationToken(record.Email())
			if err != nil {
				return apis.NewApiError(http.StatusInternalServerError, "Failed to generate verification token", nil)
			}

			// Send verification email
			verificationLink := "http://localhost:3000/confirmation/" + token
			if err := s.App.NewMailClient().Send(&mailer.Message{
				From: mail.Address{
					Address: "hello@noreply.com",
				},
				To:      []mail.Address{{Name: record.GetString("username"), Address: record.GetString("email")}},
				Subject: "Email Verification",
				Text:    "Please verify your email by clicking the link: " + verificationLink,
				HTML: fmt.Sprintf(`
					<p>Hello,</p>
					<p>Thank you for joining us at WellnessWave.</p>
					<p>Click on the button below to verify your email address.</p>
					<p>
					<a class="btn" href="%s" target="_blank" rel="noopener">Verify</a>
					</p>
					<p>
					Thanks,<br/>
					WellnessWave team
					</p>
				`, verificationLink),
			},
			); err != nil {
				return apis.NewApiError(http.StatusInternalServerError, fmt.Sprintf("Failed to send email:%s", err.Error()), err)
			}
		}
		return nil
	})
}

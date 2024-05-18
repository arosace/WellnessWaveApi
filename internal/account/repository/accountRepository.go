package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/models"
)

// UserRepository defines the interface for user data access.
type AccountRepository interface {
	Add(ctx echo.Context, user model.Account) (*models.Record, error)
	Update(ctx echo.Context, user *model.Account) (*model.Account, error)
	UpdateAuth(ctx echo.Context, user *model.Account) (*model.Account, error)
	UpdateVerify(echo.Context, *models.Record) error
	List(ctx echo.Context) ([]*model.Account, error)
	FindByID(ctx echo.Context, id string) (*model.Account, error)
	FindByEmail(ctx echo.Context, email string) (*models.Record, error)
	FindByParentID(ctx echo.Context, parentId string) ([]*model.Account, error)
}

type AccountRepo struct {
	App *pocketbase.PocketBase
}

func NewAccountRepository(app *pocketbase.PocketBase) *AccountRepo {
	return &AccountRepo{App: app}
}

func (r *AccountRepo) Add(ctx echo.Context, account model.Account) (*models.Record, error) {
	collection, err := r.App.Dao().FindCollectionByNameOrId("accounts")
	if err != nil {
		return nil, err
	}

	record := models.NewRecord(collection)
	r.LoadFromAccount(record, &account)

	if err := r.App.Dao().SaveRecord(record); err != nil {
		return nil, fmt.Errorf("Failed to save user: %w", err)
	}

	return record, nil
}

func (r *AccountRepo) List(ctx echo.Context) ([]*model.Account, error) {
	/*r.mux.Lock()
	defer r.mux.Unlock()

	var accounts []*model.Account
	for _, a := range r.accounts {
		accounts = append(accounts, a)
	}

	return accounts, nil*/
	return nil, nil
}

// FindByID returns a user by their ID.
func (r *AccountRepo) FindByID(ctx echo.Context, id string) (*model.Account, error) {
	/*r.mux.RLock()
	defer r.mux.RUnlock()

	var user *model.Account
	for _, a := range r.accounts {
		if a.ID == id {
			user = a
			return user, nil
		}
	}

	return nil, errors.New("not_found")*/
	return nil, nil
}

// FindByEmail returns a user by their email.
func (r *AccountRepo) FindByEmail(ctx echo.Context, email string) (*models.Record, error) {
	record, err := r.App.Dao().FindAuthRecordByEmail("accounts", email)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve record by email: %w", err)
	}

	if record.BaseModel.Id == "" {
		return nil, errors.New("not_found")
	}
	return record, nil
}

func (r *AccountRepo) Update(ctx echo.Context, user *model.Account) (*model.Account, error) {
	/*r.mux.Lock()
	defer r.mux.Unlock()
	r.accounts[user.Email] = user
	return user, nil*/
	return nil, nil
}

func (r *AccountRepo) UpdateAuth(ctx echo.Context, user *model.Account) (*model.Account, error) {
	/*r.mux.Lock()
	defer r.mux.Unlock()
	if _, emailHasNotChanged := r.accounts[user.Email]; emailHasNotChanged {
		r.accounts[user.Email] = user
	} else {
		for _, a := range r.accounts {
			if a.ID == user.ID {
				delete(r.accounts, a.Email)
			}
		}

		r.accounts[user.Email] = user
	}

	return user, nil*/
	return nil, nil
}

func (r *AccountRepo) UpdateVerify(ctx echo.Context, account *models.Record) error {
	account.Set("verified", true)

	if err := r.App.Dao().SaveRecord(account); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepo) FindByParentID(ctx echo.Context, parentId string) ([]*model.Account, error) {
	/*r.mux.RLock()
	defer r.mux.RUnlock()

	accounts := make([]*model.Account, 0)
	for _, a := range r.accounts {
		if a.ParentID == parentId {
			accounts = append(accounts, a)
		}
	}

	return accounts, nil*/
	return nil, nil
}

func (r *AccountRepo) LoadFromAccount(record *models.Record, account *model.Account) error {
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}

	// Convert JSON to map
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return err
	}

	record.Load(result)
	return nil
}

package repository

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/arosace/WellnessWaveApi/internal/account/model"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

// UserRepository defines the interface for user data access.
type AccountRepository interface {
	Add(ctx echo.Context, user model.Account) (*models.Record, error)
	Update(echo.Context, *models.Record) (*models.Record, error)
	Attach(echo.Context, *models.Record, string) (*models.Record, error)
	UpdateVerify(echo.Context, *models.Record) error
	List(ctx echo.Context) ([]*models.Record, error)
	FindByID(ctx echo.Context, id string) (*models.Record, error)
	FindByEmail(ctx echo.Context, email string) (*models.Record, error)
	FindByParentID(ctx echo.Context, parentId string) ([]*model.Account, error)
}

type AccountRepo struct {
	Dao *daos.Dao
}

func NewAccountRepository(dao *daos.Dao) *AccountRepo {
	return &AccountRepo{Dao: dao}
}

func (r *AccountRepo) Add(ctx echo.Context, account model.Account) (*models.Record, error) {
	collection, err := r.Dao.FindCollectionByNameOrId("accounts")
	if err != nil {
		return nil, err
	}

	record := models.NewRecord(collection)
	r.LoadFromAccount(record, &account)

	if err := r.Dao.SaveRecord(record); err != nil {
		return nil, fmt.Errorf("Failed to save user: %w", err)
	}

	return record, nil
}

func (r *AccountRepo) List(ctx echo.Context) ([]*models.Record, error) {
	query := r.Dao.RecordQuery("accounts")

	records := []*models.Record{}
	if err := query.All(&records); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *AccountRepo) Attach(ctx echo.Context, account *models.Record, parentId string) (*models.Record, error) {
	account.Set("parent_id", parentId)

	if err := r.Dao.SaveRecord(account); err != nil {
		return nil, err
	}
	return account, nil
}

// FindByID returns a user by their ID.
func (r *AccountRepo) FindByID(ctx echo.Context, id string) (*models.Record, error) {
	record, err := r.Dao.FindRecordById("accounts", id)
	if err != nil {
		return nil, err
	}
	if record.Id == "" {
		return nil, errors.New("not_found")
	}
	return record, nil
}

// FindByEmail returns a user by their email.
func (r *AccountRepo) FindByEmail(ctx echo.Context, email string) (*models.Record, error) {
	record, err := r.Dao.FindAuthRecordByEmail("accounts", email)
	if err != nil {
		return nil, fmt.Errorf("Could not retrieve record by email: %w", err)
	}

	if record.BaseModel.Id == "" {
		return nil, errors.New("not_found")
	}
	return record, nil
}

func (r *AccountRepo) Update(ctx echo.Context, account *models.Record) (*models.Record, error) {
	if err := r.Dao.SaveRecord(account); err != nil {
		return nil, err
	}
	return account, nil
}

func (r *AccountRepo) UpdateVerify(ctx echo.Context, account *models.Record) error {
	account.Set("verified", true)

	if err := r.Dao.SaveRecord(account); err != nil {
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

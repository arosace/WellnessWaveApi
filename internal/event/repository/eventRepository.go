package repository

import (
	"encoding/json"
	"fmt"

	"github.com/arosace/WellnessWaveApi/internal/event/domain"
	"github.com/arosace/WellnessWaveApi/internal/event/model"
	"github.com/labstack/echo/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
)

type EventRepo struct {
	Dao *daos.Dao
}

type EventRepository interface {
	Add(ctx echo.Context, event model.Event) (*models.Record, error)
	GetByHealthSpecialistId(echo.Context, string, string) ([]*models.Record, error)
	GetByPatientId(echo.Context, string, string) ([]*models.Record, error)
	Update(echo.Context, *models.Record) (*models.Record, error)
	GetById(echo.Context, string) (*models.Record, error)
}

func NewEventRepository(dao *daos.Dao) *EventRepo {
	return &EventRepo{
		Dao: dao,
	}
}

func (r *EventRepo) Add(ctx echo.Context, event model.Event) (*models.Record, error) {
	collection, err := r.Dao.FindCollectionByNameOrId(domain.TABLENAME)
	if err != nil {
		return nil, err
	}

	record := models.NewRecord(collection)
	r.LoadFromEvent(record, &event)
	if err := r.Dao.SaveRecord(record); err != nil {
		return nil, fmt.Errorf("Failed to save event: %w", err)
	}

	return record, nil
}

func (r *EventRepo) GetByHealthSpecialistId(ctx echo.Context, healthSpecialistId string, after string) ([]*models.Record, error) {
	column := "health_specialist_id"
	params := dbx.Params{column: healthSpecialistId}
	filter := fmt.Sprintf("%s = {:%s}", column, column)
	if after != "" {
		filter += " && event_date > {:after}"
		params["after"] = after
	}

	records, err := r.Dao.FindRecordsByFilter(
		domain.TABLENAME,
		filter,
		"-event_date",
		-1,
		0,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving events by %s: %w", column, err)
	}

	return records, nil
}

func (r *EventRepo) GetByPatientId(ctx echo.Context, patientId string, after string) ([]*models.Record, error) {
	column := "patient_id"
	params := dbx.Params{column: patientId}
	filter := fmt.Sprintf("%s = {:%s}", column, column)
	if after != "" {
		filter += " && event_date > {:after}"
		params["after"] = after
	}

	records, err := r.Dao.FindRecordsByFilter(
		domain.TABLENAME,
		filter,
		"-event_date",
		-1,
		0,
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving events by %s: %w", column, err)
	}

	return records, nil
}

func (r *EventRepo) GetById(ctx echo.Context, id string) (*models.Record, error) {
	record, err := r.Dao.FindRecordById(domain.TABLENAME, id)
	if err != nil {
		return nil, fmt.Errorf("there was an error retrieving event [%s]: %w", id, err)
	}
	return record, nil
}

func (r *EventRepo) Update(ctx echo.Context, record *models.Record) (*models.Record, error) {
	record.MarkAsNotNew()
	if err := r.Dao.SaveRecord(record); err != nil {
		return nil, fmt.Errorf("there was an error rescheduling event: %w", err)
	}
	return record, nil
}

func (r *EventRepo) LoadFromEvent(record *models.Record, event *model.Event) error {
	data, err := json.Marshal(event)
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

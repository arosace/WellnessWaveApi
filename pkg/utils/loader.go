package utils

import (
	"encoding/json"

	"github.com/pocketbase/pocketbase/models"
)

func LoadFromStruct(record *models.Record, obj interface{}) error {
	data, err := json.Marshal(obj)
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

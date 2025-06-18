package writer

import (
	"encoding/json"
	"os"

	"github.com/goutham80808/valiform/internal/validator"
)

func WriteJSON(filePath string, records []validator.Record) error {
	jsonData, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	// 0644 is a standard file permission (readable by everyone, writable only by the owner).
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}

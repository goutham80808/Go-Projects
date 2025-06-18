package reader

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/goutham80808/valiform/internal/validator"
)

// ReadCSV reads a CSV file and converts its rows into a slice of generic Records.
func ReadCSV(filePath string, hasHeader bool) ([]validator.Record, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []validator.Record{}, nil
	}

	var header []string
	var dataRows [][]string

	if hasHeader {
		header = data[0]
		dataRows = data[1:]
	} else {
		header = make([]string, len(data[0]))
		for i := range data[0] {
			header[i] = fmt.Sprintf("column_%d", i+1)
		}
		dataRows = data
	}

	var records []validator.Record
	for _, row := range dataRows {
		record := make(validator.Record)
		for i, cell := range row {
			if i < len(header) {
				record[header[i]] = cell
			}
		}
		records = append(records, record)
	}

	return records, nil
}

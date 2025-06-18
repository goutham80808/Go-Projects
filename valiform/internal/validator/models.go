package validator

import "fmt"

// Record is a generic representation of a single row of data from any file type.
// We use a map where the key is the field name (e.g., "StudentID", "Grade")
// and the value is the data for that field.
type Record map[string]interface{}

// ValidationError stores the details of a single failed validation.
type ValidationError struct {
	RowNumber int    // The line number in the file where the error occurred
	FieldName string // The name of the field that failed validation
	Message   string // A human-readable message explaining the failure
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("row %d, field '%s': %s", e.RowNumber, e.FieldName, e.Message)
}

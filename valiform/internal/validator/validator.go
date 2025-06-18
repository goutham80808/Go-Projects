package validator

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/goutham80808/valiform/internal/config"
)

// ValidateRecord checks a single record against a set of field rules.
// It returns a slice of ValidationErrors. If the slice is empty, the record is valid.
func ValidateRecord(record Record, fieldRules []config.FieldRule, rowNum int) []ValidationError {
	var errors []ValidationError

	for _, rule := range fieldRules {
		value, ok := record[rule.Name]

		// 1. Check for 'required'
		if rule.Rules.Required != nil && *rule.Rules.Required {
			if !ok || value == "" {
				errors = append(errors, ValidationError{
					RowNumber: rowNum,
					FieldName: rule.Name,
					Message:   "field is required but is missing or empty",
				})
				// If a required field is missing, no other checks can be done on it.
				continue
			}
		}

		// If the field is not present and not required, we can skip other checks.
		if !ok {
			continue
		}

		valStr := fmt.Sprintf("%v", value)

		// 2. Check by 'type' and run type-specific rules
		switch rule.Type {
		case "integer":
			i, err := strconv.Atoi(valStr)
			if err != nil {
				errors = append(errors, ValidationError{rowNum, rule.Name, "value is not a valid integer"})
				continue // Can't do min/max if it's not an int
			}
			// Check for 'min'
			if rule.Rules.Min != nil && i < *rule.Rules.Min {
				errors = append(errors, ValidationError{rowNum, rule.Name, fmt.Sprintf("value %d is less than min %d", i, *rule.Rules.Min)})
			}
			// Check for 'max'
			if rule.Rules.Max != nil && i > *rule.Rules.Max {
				errors = append(errors, ValidationError{rowNum, rule.Name, fmt.Sprintf("value %d is greater than max %d", i, *rule.Rules.Max)})
			}

		case "string":
			// 3. Check for 'regex'
			if rule.Rules.Regex != nil {
				re := regexp.MustCompile(*rule.Rules.Regex)
				if !re.MatchString(valStr) {
					errors = append(errors, ValidationError{rowNum, rule.Name, fmt.Sprintf("value '%s' does not match regex pattern '%s'", valStr, *rule.Rules.Regex)})
				}
			}
		}
	}

	return errors
}

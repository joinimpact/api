package validators

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

// MinAge validates a date by a minimum age (in years) relative to the
// current date.
func MinAge(fl validator.FieldLevel) bool {
	// Convert the tag to an integer.
	tag := fl.Param()
	fmt.Println(tag)
	i, err := strconv.ParseInt(tag, 10, 32)
	if err != nil {
		return false
	}

	years := int(i)

	// Get the time.Time from the field.
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Create a comparison date by subtracting the specified years.
	compareDate := time.Now().AddDate(-1*years, 0, 0)

	// Return true if the date is before or equal to the comparison date.
	return date.Before(compareDate) || date.Equal(compareDate)
}

// MaxAge validates a date by a maximum age (in years) relative to the
// current date.
func MaxAge(fl validator.FieldLevel) bool {
	// Convert the tag to an integer.
	tag := fl.Param()
	i, err := strconv.ParseInt(tag, 10, 32)
	if err != nil {
		return false
	}

	years := int(i)

	// Get the time.Time from the field.
	date, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	// Create a comparison date by subtracting the specified years.
	compareDate := time.Now().AddDate(years, 0, 0)

	// Return true if the date is before or equal to the comparison date.
	return date.Before(compareDate) || date.Equal(compareDate)
}

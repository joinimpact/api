package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

const urlRegex = "^(http:\\/\\/www\\.|https:\\/\\/www\\.|http:\\/\\/|https:\\/\\/)?[a-z0-9]+([\\-\\.]{1}[a-z0-9]+)*\\.[a-z]{2,5}(:[0-9]{1,5})?(\\/.*)?$/"

// URL validates a web URL.
func URL(fl validator.FieldLevel) bool {
	if match, err := regexp.MatchString(urlRegex, fl.Field().Interface().(string)); !match || err != nil {
		return false
	}

	return true
}

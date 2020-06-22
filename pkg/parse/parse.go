// Package parse exports methods for parsing and validating the bodies of HTTP requests.
package parse

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/joinimpact/api/pkg/resp"
	"github.com/mitchellh/mapstructure"
	validator "gopkg.in/go-playground/validator.v9"
)

// POST parses and reflects data into a struct.
func POST(w http.ResponseWriter, r *http.Request, s interface{}) error {
	ct := r.Header.Get("Content-Type")
	ct = strings.Split(ct, ";")[0]
	if ct == "multipart/form-data" {
		o := mapstructure.DecoderConfig{
			TagName:          "json",
			WeaklyTypedInput: true,
			Result:           s,
		}

		r.ParseMultipartForm(4096)
		form := map[string]interface{}{}
		for key, value := range r.Form {
			form[key] = value[0]
		}
		decoder, err := mapstructure.NewDecoder(&o)
		if err != nil {
			log.Println(err)
		}

		err = decoder.Decode(form)
		if err != nil {
			log.Printf("An error occurred while parsing form data: %e", err)
			resp.BadRequest(w, r, resp.Err{
				Code:    100,
				Message: "unable to decode input",
			})
			return err
		}
	} else if ct == "application/json" {
		err := json.NewDecoder(r.Body).Decode(s)
		if err != nil {
			resp.BadRequest(w, r, resp.Err{
				Code:    100,
				Message: "unable to decode input",
			})
			return err
		}
	}

	validate := validator.New()
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	fieldErr := resp.Err{}

	fieldErr.Code = 98
	fieldErr.Message = "missing or invalid fields"
	fieldErr.InvalidFields = []string{}

	for _, err := range err.(validator.ValidationErrors) {
		field, ok := reflect.TypeOf(s).Elem().FieldByName(err.StructField())
		if ok {
			fieldErr.InvalidFields = append(fieldErr.InvalidFields, field.Tag.Get("json"))
		} else {
			fieldErr.InvalidFields = append(fieldErr.InvalidFields, err.Field())
		}
	}

	if len(fieldErr.InvalidFields) > 0 {
		resp.BadRequest(w, r, fieldErr)
		return errors.New("validation error")
	}

	return nil
}

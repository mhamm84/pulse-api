package api

import (
	"encoding/json"
	"github.com/mhamm84/pulse-api/internal/validator"
	"net/http"
	"net/url"
	"strconv"
)

type envelope map[string]interface{}

func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func (app *application) WriteJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}
	js = append(js, '\n')

	for k, v := range headers {
		w.Header()[k] = v
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

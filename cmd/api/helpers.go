package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]interface{}

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

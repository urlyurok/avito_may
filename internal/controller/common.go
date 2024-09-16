package controller

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"io"
	"net/http"
)

func ParseJSONBody[T any](r *http.Request, w http.ResponseWriter) (*T, error) {
	var t T
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &t)
	if err != nil {
		return nil, err
	}

	var validate = validator.New()

	err = validate.Struct(t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func DecodeFormParams[T any](r *http.Request) (*T, error) {
	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	params := new(T)
	var decoder = schema.NewDecoder()
	if err := decoder.Decode(params, r.Form); err != nil {
		return nil, err
	}

	validate := validator.New()

	if err := validate.Struct(params); err != nil {
		return nil, err
	}

	return params, nil
}

func SendJSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

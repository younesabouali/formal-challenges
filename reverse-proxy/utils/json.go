package utils

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func ParseInt32(stringifiedNumber string) (int32, error) {
	params, err := strconv.ParseInt(stringifiedNumber, 10, 32)
	if err != nil {
		return 0, err
	}
	e := int32(params)
	return e, err
}
func UrlParamsParser(r *http.Request, fieldName string) string {
	return r.URL.Query().Get(fieldName)

}
func BodyParser[T interface{}](r *http.Request, e T) (T, error) {
	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		return e, errors.New("Couldn't parse params")
	}
	return e, nil

}
func RespondWithArrayJSON[T interface{}](w http.ResponseWriter, code int, payload []T) {
	if payload != nil {
		RespondWithJSON(w, code, payload)
		return
	}

	var empyFollow [0]int
	RespondWithJSON(w, code, empyFollow)
}
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responded with 500 error ", msg)
	}
	type responseError struct {
		Error string `json:"error"`
	}
	RespondWithJSON(w, code, responseError{Error: msg})

}

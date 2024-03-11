package utilities

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wrapper := make(map[string]interface{})
	wrapper[wrap] = data

	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func ErrorJSON(w http.ResponseWriter, err error) {
	type JSONError struct {
		Message string `json:"message"`
	}

	theError := JSONError{
		Message: err.Error(),
	}

	WriteJSON(w, http.StatusBadRequest, theError, "error")
}

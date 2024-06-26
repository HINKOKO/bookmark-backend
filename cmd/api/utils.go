package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"net/http"
	"regexp"
)

// JSONResponse - structure to pack the json response data
type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Character set from which to generate the random string (email validation)
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// writeJSON - write a JSON response to the application
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// readJSON - Read JSON from the application
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}

// errorJSON - writes an error in JSON format to be send throughout application
func (app *application) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload JSONResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}

// generateRandomString - generate a random string for new user wishing to register
func generateRandomString(length int) string {
	randomString := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range randomString {
		randowIdx, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return ""
		}
		randomString[i] = charset[randowIdx.Int64()]
	}
	return string(randomString)
}

// isValidURL - validates an Url against regular expression
func isValidURL(url string) bool {
	re := regexp.MustCompile(`^(https?://)?((([a-z\d]([a-z\d-]*[a-z\d])*)\.?)+[a-z]{2,}|(\d{1,3}\.){3}\d{1,3})(:\d+)?(/[-a-z\d%_.~+]*)*(\?[;&a-z\d%_.~+=-]*)?(#[-a-z\d_]*)?$`)
	return re.MatchString(url)
}

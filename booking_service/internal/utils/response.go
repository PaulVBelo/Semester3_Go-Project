package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// JSONResponse отправляет успешный JSON-ответ
func JSONResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

// ErrorResponse отправляет JSON-ответ с ошибкой
func ErrorResponse(w http.ResponseWriter, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, jsonError(message, err.Error()), http.StatusInternalServerError)
	} else {
		http.Error(w, jsonError(message, ""), http.StatusBadRequest)
	}
}

func jsonError(message, details string) string {
	response := map[string]string{"error": message}
	if details != "" {
		response["details"] = details
	}
	output, _ := json.Marshal(response)
	return string(output)
}

func ParseJSONBody(r *http.Request, dest interface{}) error {
	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read request body: " + err.Error())
	}
	defer r.Body.Close()

	// Декодирование JSON в структуру
	if err := json.Unmarshal(body, dest); err != nil {
		return errors.New("failed to unmarshal JSON: " + err.Error())
	}

	return nil
}

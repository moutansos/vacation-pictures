package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

type ErrorData struct {
	Level string
	Message string
}

func ErrorHandler(logger *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postErrors(w, r, logger)
			break
		default:
            logger.Info("Method not allowed", "method", r.Method, "ip", r.RemoteAddr)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func postErrors(w http.ResponseWriter, r *http.Request, logger *slog.Logger) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Error reading body", "error", err)
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	errorData := []ErrorData{}
	err = json.Unmarshal(body, &errorData)

	if err != nil {
		logger.Error("Error unmarshalling body", "error", err)
		http.Error(w, "Error unmarshalling body", http.StatusBadRequest)
		return
	}

	for _, e := range errorData {
		errorMessage := "Frontend Event:\n" + e.Message
		switch e.Level {
		case "error":
			logger.Error(errorMessage)
			break
		case "warn":
			logger.Warn(errorMessage)
			break
		case "info":
			logger.Info(errorMessage)
			break
		default:
			logger.Error("Unknown log level", "level", e.Level)
		}
	}

	w.WriteHeader(http.StatusOK)
}

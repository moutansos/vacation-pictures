package handlers

import (
	"fmt"
	"net/http"
    "log/slog"
)
const isDevMode = true;

func logError(w http.ResponseWriter, err error, msg string, logger *slog.Logger) {
    message := "An error occured"
    if(msg != "") {
        message = msg
    }

    if(isDevMode) {
        errorMsg := fmt.Sprintf("%s: %s", message, err)
        logger.Error("An error occured", "error", err)
        http.Error(w, errorMsg, http.StatusInternalServerError)
    } else {
        fullErrorMsg := fmt.Sprintf("%s: %s", message, err)
        publicErrorMsg := message

        logger.Error("An error occured", "error", fullErrorMsg)
        http.Error(w, publicErrorMsg, http.StatusInternalServerError)
    }
}

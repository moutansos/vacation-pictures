package handlers

import (
	"fmt"
	"log"
	"net/http"
)

const isDevMode = true;

func logError(w http.ResponseWriter, err error, msg string) {
    message := "An error occured"
    if(msg != "") {
        message = msg
    }

    if(isDevMode) {
        errorMsg := fmt.Sprintf("%s: %s", message, err)
        log.Printf("%s", errorMsg)
        http.Error(w, errorMsg, http.StatusInternalServerError)
    } else {
        fullErrorMsg := fmt.Sprintf("%s: %s", message, err)
        publicErrorMsg := message

        log.Printf("%s", fullErrorMsg)
        http.Error(w, publicErrorMsg, http.StatusInternalServerError)
    }
}

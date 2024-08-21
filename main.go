package main

import (
	"log"
	"mime"
	"net/http"
	"vacation-pictures/handlers"
	"vacation-pictures/infra"
)

func main() {
    fixMimeTypes()
    log.Println("Starting vacations application...")
    db, err := infra.ConnectDb("vacations.json")
    if err != nil {
        log.Fatalf("Error loading db: %s", err)
    }

    vacations, err := db.GetVacations()
    if err != nil {
        log.Fatalf("Error loading vacations from db: %s", err)
    }

    vacaCount := len(vacations)
    log.Printf("Found %d vacations in db\n", vacaCount)

    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    http.HandleFunc("/", handlers.IndexHandler(db))
    http.HandleFunc("/vacations", handlers.VacationHandler(db))

    log.Fatal(http.ListenAndServe(":8080", nil))
}

func fixMimeTypes() {
    err1 := mime.AddExtensionType(".js", "text/javascript")
    if err1 != nil {
        log.Printf("Error in mime js %s", err1.Error())
    }

    err2 := mime.AddExtensionType(".css", "text/css")
    if err2 != nil {
        log.Printf("Error in mime js %s", err2.Error())
    }
}

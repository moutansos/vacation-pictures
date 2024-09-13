package main

import (
	"log"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"vacation-pictures/handlers"
	"vacation-pictures/infra"

	charmlog "github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	slogmulti "github.com/samber/slog-multi"
	slogslack "github.com/samber/slog-slack/v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
        log.Println("No .env file found. No environment variables loaded from file")
	}

    var loggerWebhookUrl = os.Getenv("SLACK_WEBHOOK_URL")
	const channel = "alerts"
	logger := slog.New(slogmulti.Fanout(
		slogslack.Option{Level: slog.LevelInfo, WebhookURL: loggerWebhookUrl, Channel: channel, AddSource: true}.NewSlackHandler(),
		charmlog.NewWithOptions(os.Stdout, charmlog.Options{ReportCaller: true, ReportTimestamp: true}),
	))
	logger = logger.
		With("app", "vacation-pictures")

	logger.Info("Starting vacations application...")
	fixMimeTypes(logger)

	db, err := infra.ConnectDb("vacations.json")
	if err != nil {
		logger.Error("Error loading db", "error", err)
		panic(err)
	}

	vacations, err := db.GetVacations()
	if err != nil {
		logger.Error("Error loading vacations from db", "error", err)
		panic(err)
	}

	vacaCount := len(vacations)
	logger.Info("Vacations Retrieved", "vacationCount", vacaCount)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handlers.IndexHandler(db, logger))
	http.HandleFunc("/vacations", handlers.VacationHandler(db, logger))

	err = http.ListenAndServe(":8081", nil)
	if err != nil {
		logger.Error("Error starting server", "error", err)
		panic(err)
	}
}

func fixMimeTypes(loger *slog.Logger) {
	err1 := mime.AddExtensionType(".js", "text/javascript")
	if err1 != nil {
		loger.Error("Error adding mime type js", "error", err1)
	}

	err2 := mime.AddExtensionType(".css", "text/css")
	if err2 != nil {
		loger.Error("Error in mime type css", "error", err2)
	}
}

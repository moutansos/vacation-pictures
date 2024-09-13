package handlers

import (
	"bufio"
	"bytes"
	"net/http"
	"log/slog"
	"vacation-pictures/data"
	"vacation-pictures/infra"

	"github.com/philippta/go-template/html/template"
)

func IndexHandler(db *infra.Db, logger *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
            getIndex(db, w, logger)
			break
		default:
            logger.Info("Method not allowed", "method", r.Method, "ip", r.RemoteAddr)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

type IndexPageData struct {
    Vacations []data.Vacation
}

func getIndex(db *infra.Db, w http.ResponseWriter, logger *slog.Logger) {
    tmpl, err := template.ParseFiles(
        "pages/index.html",
    )
    if err != nil {
        logError(w, err, "Error parsing templates", logger)
        return
    }

    vacations, err := db.GetVacations()
    if err != nil {
        logError(w, err, "Error retriving vacations from db", logger)
        return
    }

    pageData := IndexPageData {
        Vacations: vacations,
    }

    var b bytes.Buffer
    templateBuff := bufio.NewWriter(&b)
    err = tmpl.Execute(templateBuff, pageData)
    if err != nil {
        logError(w, err, "Unable to execute template for this page", logger)
        return
    }

    err = templateBuff.Flush()
    if err != nil {
        logError(w, err, "", logger)
    }

    w.Header().Set("Content-Type", "text/html")
    w.Write(b.Bytes())
}

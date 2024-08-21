package handlers

import (
	"bufio"
	"bytes"
	"net/http"
	"vacation-pictures/data"
	"vacation-pictures/infra"

	"github.com/philippta/go-template/html/template"
)

func IndexHandler(db *infra.Db) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
            getIndex(db, w)
			break
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

type IndexPageData struct {
    Vacations []data.Vacation
}

func getIndex(db *infra.Db, w http.ResponseWriter) {
    tmpl, err := template.ParseFiles(
        "pages/index.html",
    )
    if err != nil {
        logError(w, err, "Error parsing templates")
        return
    }

    vacations, err := db.GetVacations()
    if err != nil {
        logError(w, err, "Error retriving vacations from db")
        return
    }

    pageData := IndexPageData {
        Vacations: vacations,
    }

    var b bytes.Buffer
    templateBuff := bufio.NewWriter(&b)
    err = tmpl.Execute(templateBuff, pageData)
    if err != nil {
        logError(w, err, "Unable to execute template for this page")
        return
    }

    err = templateBuff.Flush()
    if err != nil {
        logError(w, err, "")
    }

    w.Header().Set("Content-Type", "text/html")
    w.Write(b.Bytes())
}

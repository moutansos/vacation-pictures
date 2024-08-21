package handlers

import (
	"bufio"
	"bytes"
	"net/http"
	"vacation-pictures/data"
	"vacation-pictures/infra"

	"github.com/philippta/go-template/text/template"
)

func VacationHandler(db *infra.Db) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getVacation(db, w, r)
            break
        default:
            http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        }
    }
}

type VacationPageData struct {
    Vacation *data.Vacation
}

func getVacation(db *infra.Db, w http.ResponseWriter, r *http.Request) {
    tmpl, err := template.ParseFiles(
        "pages/vacation.html",
    )

    if err != nil {
        logError(w, err, "Error parsing templates")
        return
    }

    values := r.URL.Query()
    id := values.Get("id")

    if id == "" {
        http.Error(w, "Id was blank!", http.StatusNotFound)
        return
    }


    vacation, err := db.GetVacationById(id)
    if err != nil {
        http.Error(w, "Vacation not found", http.StatusNotFound)
        return
    }

    pageData := VacationPageData {
        Vacation: vacation,
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
        return
    }

    w.Header().Set("Content-Type", "text/html")
    w.Write(b.Bytes())
}

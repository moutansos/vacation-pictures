package handlers

import (
	"bufio"
	"bytes"
	"net/http"
	"strconv"
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
    CurrentPic *data.VacaPicture
    NextPicIndex *int
    PrevPicIndex *int
    CurrentPicIndex int
	CurrentPicStyle string
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
    pic := values.Get("pic")

    if id == "" {
        http.Error(w, "Id was blank!", http.StatusNotFound)
        return
    }

    if pic == "" {
        pic = "0"
    }


    vacation, err := db.GetVacationById(id)
    if err != nil {
        http.Error(w, "Vacation not found", http.StatusNotFound)
        return
    }

    picIndex, err := strconv.Atoi(pic)
    if err != nil {
        http.Error(w, "Invalid picture index", http.StatusBadRequest)
        return
    }
    

    if picIndex < 0 || picIndex > len(vacation.Pictures) {
        http.Error(w, "Picture index out of range", http.StatusBadRequest)
        return
    }

    currentPic := vacation.Pictures[picIndex]

    var nextPicIndex *int
    var prevPicIndex *int

    calculatedNextPicIndex := picIndex + 1
    calculatedPrevPicIndex := picIndex - 1

    if calculatedNextPicIndex < len(vacation.Pictures) {
        nextPicIndex = &calculatedNextPicIndex
    }

    if calculatedPrevPicIndex >= 0 {
        prevPicIndex = &calculatedPrevPicIndex
    }

    pageData := VacationPageData {
        Vacation: vacation,
        CurrentPic: &currentPic,
        NextPicIndex: nextPicIndex,
        PrevPicIndex: prevPicIndex,
        CurrentPicIndex: picIndex,
		CurrentPicStyle: "",
    }

	if currentPic.Rotate != nil {
		pageData.CurrentPicStyle = "transform: rotate(" + *currentPic.Rotate + ");"
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

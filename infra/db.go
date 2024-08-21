package infra

import (
	"encoding/json"
	"fmt"
	"os"
    "vacation-pictures/data"
)

type Db struct {
    VacationsRoot data.VacationsRoot `json:"vacationsRoot"`
}

func ConnectDb(filePath string) (*Db, error) {
    var jsonRoot data.VacationsRoot
    configFile, err := os.Open(filePath)
    if err != nil {
        return nil, err
    }

    jsonParser := json.NewDecoder(configFile)
    if err = jsonParser.Decode(&jsonRoot); err != nil {
        return nil, err
    }

    db := Db{ VacationsRoot: jsonRoot }
    return &db, nil
}

func (db *Db) GetVacations() ([]data.Vacation, error) {
    return db.VacationsRoot.Vacations, nil
}

func (db *Db) GetVacationById(id string) (*data.Vacation, error) {
    for _, vaca := range db.VacationsRoot.Vacations {
        if vaca.ID == id {
            return &vaca, nil
        }
    }

    return nil, fmt.Errorf("No vacation with ID: '%s' was found", id)
}

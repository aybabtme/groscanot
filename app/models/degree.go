package models

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/db"
	"log"
	"time"
)

const DegreeCollection = "degrees"

type Degree struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Url         string    `json:"url"`
	Credit      int       `json:"credit"`
	Mandatory   []string  `json:"mandat"`
	Extra       []string  `json:"extra"`
	LastUpdated time.Time `json:"updated"`
}

func DegreeGetAll() ([]Degree, error) {
	payloads, err := db.Db.GetAll(DegreeCollection)
	if err != nil {
		log.Printf("Error Db.GetAll, %v", err)
		return nil, err
	}
	var d Degree
	var degrees []Degree
	for _, payload := range payloads {
		if err = json.Unmarshal(payload, &d); err != nil {
			log.Printf("Error unmarshalling, %v", err)
			return degrees, err
		}
		degrees = append(degrees, d)
	}
	return degrees, nil

}

func DegreeGetJson(degreeName string) (string, bool, error) {
	payload, ok, err := db.Db.Get(DegreeCollection + db.KeySep + degreeName)
	if err != nil || !ok {
		return "", ok, err
	}

	return string(payload), true, err
}

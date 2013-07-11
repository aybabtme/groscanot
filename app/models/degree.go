package models

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/db"
	"time"
)

const DegreeCollection = "degrees"

type Degree struct {
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
		return nil, err
	}
	var d *Degree
	var degrees []Degree
	for _, payload := range payloads {
		if err = json.Unmarshal(payload, d); err != nil {
			return degrees, err
		}
		degrees = append(degrees, *d)
	}
	return degrees, nil

}

func DegreeGet(degreeName string) (Degree, bool, error) {
	var d Degree
	payload, ok, err := db.Db.Get(DegreeCollection + db.KeySep + degreeName)
	if err != nil || !ok {
		return d, ok, err
	}

	if err = json.Unmarshal(payload, &d); err != nil {
		return d, false, err
	}
	return d, true, err
}

package models

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/db"
	"time"
)

const CourseCollection = "courses"

type Course struct {
	Id          string    `json:"id"`
	Topic       string    `json:"topic"`
	Code        string    `json:"code"`
	Url         string    `json:"url"`
	Level       int       `json:"level"`
	Credit      int       `json:"credit"`
	Name        string    `json:"name"`
	Description string    `json:"descr"`
	Dependency  []string  `json:"depend"`
	Equivalence []string  `json:"equiv"`
	LastUpdated time.Time `json:"updated"`
}

func CourseGetAll() ([]Course, error) {
	payloads, err := db.Db.GetAll(CourseCollection)
	if err != nil {
		return nil, err
	}
	var c *Course
	var courses []Course
	for _, payload := range payloads {
		if err = json.Unmarshal(payload, c); err != nil {
			return courses, err
		}
		courses = append(courses, *c)
	}
	return courses, nil

}

func CourseGet(courseCode string) (Course, bool, error) {
	var c Course
	payload, ok, err := db.Db.Get(CourseCollection + db.KeySep + courseCode)
	if err != nil || !ok {
		return c, ok, err
	}

	if err = json.Unmarshal(payload, &c); err != nil {
		return c, false, err
	}
	return c, true, err
}

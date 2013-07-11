package models

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/db"
	"log"
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
		log.Printf("Error Db.GetAll, %v", err)
		return nil, err
	}
	var c Course
	var courses []Course
	for _, payload := range payloads {
		if err = json.Unmarshal(payload, &c); err != nil {
			log.Printf("Error unmarshalling, %v", err)
			return courses, err
		}
		courses = append(courses, c)
	}
	return courses, nil

}

func CourseGetJson(courseCode string) (string, bool, error) {
	fullkey := CourseCollection + db.KeySep + courseCode
	log.Printf("Getting %s", fullkey)
	payload, ok, err := db.Db.Get(fullkey)
	if err != nil || !ok {
		log.Printf("!!! ok=%v err=%v", ok, err)
		return "", ok, err
	}

	return string(payload), true, err
}

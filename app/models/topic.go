package models

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/db"
	"log"
	"time"
)

const TopicCollection = "topics"

type TopicShort struct {
	Code        string `json:"code"`
	Description string `json:"descr"`
}

type Topic struct {
	*TopicShort
	Courses     []string  `json:"courses"`
	LastUpdated time.Time `json:"updated"`
}

func TopicGetAll() ([]TopicShort, error) {
	payloads, err := db.Db.GetAll(TopicCollection)
	if err != nil {
		log.Printf("Error Db.GetAll, %v", err)
		return nil, err
	}
	var t TopicShort
	var topics []TopicShort
	for _, payload := range payloads {
		if err = json.Unmarshal(payload, &t); err != nil {
			log.Printf("Error unmarshalling, %v", err)
			return topics, err
		}
		topics = append(topics, t)
	}
	return topics, nil

}

func TopicGetJson(topicCode string) (string, bool, error) {
	payload, ok, err := db.Db.Get(TopicCollection + db.KeySep + topicCode)
	if err != nil || !ok {
		return "", ok, err
	}
	return string(payload), true, err
}

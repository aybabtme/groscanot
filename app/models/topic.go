package models

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/db"
	"time"
)

const TopicCollection = "topic"

type Topic struct {
	Code        string    `json:"code"`
	Description string    `json:"descr"`
	LastUpdated time.Time `json:"updated"`
}

func TopicGetAll() ([]Topic, error) {
	payloads, err := db.Db.GetAll(TopicCollection)
	if err != nil {
		return nil, err
	}
	var t *Topic
	var topics []Topic
	for _, payload := range payloads {
		if err = json.Unmarshal(payload, t); err != nil {
			return topics, err
		}
		topics = append(topics, *t)
	}
	return topics, nil

}

func TopicGet(topicCode string) (Topic, bool, error) {
	var t Topic
	payload, ok, err := db.Db.Get(TopicCollection + db.KeySep + topicCode)
	if err != nil || !ok {
		return t, ok, err
	}

	if err = json.Unmarshal(payload, &t); err != nil {
		return t, false, err
	}
	return t, true, err
}

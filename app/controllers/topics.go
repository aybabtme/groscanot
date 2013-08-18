package controllers

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/models"
	"github.com/robfig/revel"
	"time"
)

type Topics struct {
	*revel.Controller
}

func (t *Topics) Index() revel.Result {
	defer gatherMetrics(t.Request, time.Now())
	topics, err := models.TopicGetAll()
	if err != nil {
		revel.ERROR.Printf("Error from models.TopicGetAll: %v", err)
		return t.Forbidden("This resources is not available to you")
	}

	setJSONMimeType(t.Controller)

	jsonpRequest := t.Request.URL.Query().Get("callback")
	if jsonpRequest == "" {
		return t.RenderJson(topics)
	}

	rawTopics, err := json.Marshal(topics)
	if err != nil {
		revel.ERROR.Printf("Error marshalling topics, %v", err)
		return t.RenderError(err)
	}
	return t.RenderText(jsonpRequest + "(" + string(rawTopics) + ")")
}

func (t *Topics) Get(code string) revel.Result {
	defer gatherMetrics(t.Request, time.Now())
	topic, ok, err := models.TopicGetJson(code)
	if err != nil {
		revel.ERROR.Printf("Error from models.TopicGetJson, %v", err)
		return t.Forbidden("This resource is not available to you: %v", code)
	}
	if !ok {
		return t.NotFound("This topic is unknown: %v", code)
	}

	setJSONMimeType(t.Controller)

	jsonpRequest := t.Request.URL.Query().Get("callback")
	if jsonpRequest == "" {
		return t.RenderText(topic)
	}

	return t.RenderText(jsonpRequest + "(" + topic + ")")
}

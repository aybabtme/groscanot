package controllers

import (
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
	return t.RenderJson(topics)
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
	return t.RenderText(topic)
}

package controllers

import (
	"github.com/aybabtme/groscanot/app/models"
	"github.com/robfig/revel"
	"log"
	"time"
)

type Topics struct {
	*revel.Controller
}

func (t *Topics) Index() revel.Result {
	t0 := time.Now()
	topics, err := models.TopicGetAll()
	if err != nil {
		log.Printf("Error from models.TopicGetAll: %v", err)
		return t.Forbidden("This resources is not available to you")
	}
	log.Printf("Done in %v", time.Since(t0))
	return t.RenderJson(topics)
}

func (t *Topics) Get(code string) revel.Result {
	t0 := time.Now()
	topic, ok, err := models.TopicGetJson(code)
	if err != nil {
		log.Printf("Error from models.TopicGetJson, %v", err)
		return t.Forbidden("This resource is not available to you: %v", code)
	}
	if !ok {
		return t.NotFound("This topic is unknown: %v", code)
	}
	log.Printf("Done in %v", time.Since(t0))
	return t.RenderText(topic)
}

package controllers

import (
	"github.com/aybabtme/groscanot/app/models"
	"github.com/robfig/revel"
	"log"
	"time"
)

type Courses struct {
	*revel.Controller
}

func (c *Courses) Index() revel.Result {
	t0 := time.Now()
	courses, err := models.CourseGetAll()
	if err != nil {
		log.Printf("Error from models.CourseGetAll: %v", err)
		return c.Forbidden("This resource is not available to you")
	}
	log.Printf("Done in %v", time.Since(t0))
	return c.RenderJson(courses)
}

func (c *Courses) Get(code string) revel.Result {
	t0 := time.Now()
	course, ok, err := models.CourseGetJson(code)
	if err != nil {
		log.Printf("Error from models.CourseGetJson, %v", err)
		return c.Forbidden("This resource is not available to you: %v", code)
	}
	if !ok {
		return c.NotFound("This course is unknown: %v", code)
	}
	log.Printf("Done in %v", time.Since(t0))
	return c.RenderText(course)
}

package controllers

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/models"
	"github.com/robfig/revel"
	"time"
)

type Courses struct {
	*revel.Controller
}

func (c *Courses) Index() revel.Result {
	defer gatherMetrics(c.Request, time.Now())
	courses, err := models.CourseGetAll()

	if err != nil {
		revel.ERROR.Printf("Error from models.CourseGetAll: %v", err)
		return c.Forbidden("This resource is not available to you")
	}

	setJSONMimeType(c.Controller)

	jsonpRequest := c.Request.URL.Query().Get("callback")
	if jsonpRequest == "" {
		return c.RenderJson(courses)
	}

	rawCourses, err := json.Marshal(courses)
	if err != nil {
		revel.ERROR.Printf("Error marshalling courses, %v", err)
		return c.RenderError(err)
	}
	return c.RenderText(jsonpRequest + "(" + string(rawCourses) + ")")
}

func (c *Courses) Get(code string) revel.Result {
	defer gatherMetrics(c.Request, time.Now())
	course, ok, err := models.CourseGetJson(code)
	if err != nil {
		revel.ERROR.Printf("Error from models.CourseGetJson, %v", err)
		return c.Forbidden("This resource is not available to you: %v", code)
	}

	if !ok {
		return c.NotFound("This course is unknown: %v", code)
	}

	setJSONMimeType(c.Controller)

	jsonpRequest := c.Request.URL.Query().Get("callback")
	if jsonpRequest == "" {
		return c.RenderText(course)
	}

	return c.RenderText(jsonpRequest + "(" + course + ")")
}

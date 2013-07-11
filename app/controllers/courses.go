package controllers

import (
	"github.com/aybabtme/groscanot/app/models"
	"github.com/robfig/revel"
)

type Course struct {
	*revel.Controller
}

func (c *Course) Index() {
	courses, err := models.CourseGetAll()

	return c.Render(courses)
}

func (c *Course) Get(code string) {
	course, ok, err := models.CourseGet(code)

	return c.Render(course)
}

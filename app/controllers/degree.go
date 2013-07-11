package controllers

import (
	"github.com/aybabtme/groscanot/app/models"
	"github.com/robfig/revel"
)

type Degree struct {
	*revel.Controller
}

func (d *Degree) Index() {
	degrees, err := models.DegreeGetAll()

	return d.Render(degrees)
}

func (d *Degree) Get(name string) {
	degree, ok, err := models.TopicGet(name)

	return d.Render(degree)
}

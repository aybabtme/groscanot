package controllers

import (
	"github.com/aybabtme/groscanot/app/models"
	"github.com/robfig/revel"
	"time"
)

type Degrees struct {
	*revel.Controller
}

func (d *Degrees) Index() revel.Result {
	t0 := time.Now()
	degrees, err := models.DegreeGetAll()
	if err != nil {
		revel.ERROR.Printf("Error from models.DegreeGetAll: %v", err)
		return d.Forbidden("This resource is not available to you")
	}
	revel.INFO.Printf("Index - Done in %v", time.Since(t0))
	return d.RenderJson(degrees)
}

func (d *Degrees) Get(name string) revel.Result {
	t0 := time.Now()
	degree, ok, err := models.DegreeGetJson(name)
	if err != nil {
		revel.ERROR.Printf("Error from models.DegreeGetJson, %v", err)
		return d.Forbidden("This resource is not available to you: %v", name)
	}
	if !ok {
		return d.NotFound("This degree is unknown: %v", name)
	}
	revel.INFO.Printf("Get - Done %dB in %v", len(degree), time.Since(t0))
	return d.RenderText(degree)
}

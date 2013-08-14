package controllers

import (
	"encoding/json"
	"github.com/aybabtme/groscanot/app/models"
	"github.com/robfig/revel"
	"time"
)

type Degrees struct {
	*revel.Controller
}

func (d *Degrees) Index() revel.Result {
	defer gatherMetrics(d.Request, time.Now())
	degrees, err := models.DegreeGetAll()
	if err != nil {
		revel.ERROR.Printf("Error from models.DegreeGetAll: %v", err)
		return d.Forbidden("This resource is not available to you")
	}

	jsonpRequest := d.Request.URL.Query().Get("callback")
	if jsonpRequest == "" {
		return d.RenderJson(degrees)
	}

	rawDegree, err := json.Marshal(degrees)
	if err != nil {
		revel.ERROR.Printf("Error marshalling degrees, %v", err)
		return d.RenderError(err)
	}
	return d.RenderText(jsonpRequest + "(" + string(rawDegree) + ")")
}

func (d *Degrees) Get(name string) revel.Result {
	defer gatherMetrics(d.Request, time.Now())
	degree, ok, err := models.DegreeGetJson(name)
	if err != nil {
		revel.ERROR.Printf("Error from models.DegreeGetJson, %v", err)
		return d.Forbidden("This resource is not available to you: %v", name)
	}
	if !ok {
		return d.NotFound("This degree is unknown: %v", name)
	}

	jsonpRequest := d.Request.URL.Query().Get("callback")
	if jsonpRequest == "" {
		return d.RenderText(degree)
	}

	return d.RenderText(jsonpRequest + "(" + degree + ")")
}

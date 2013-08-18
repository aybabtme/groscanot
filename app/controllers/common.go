package controllers

import (
	"github.com/robfig/revel"
	"time"
)

func gatherMetrics(r *revel.Request, t0 time.Time) {
	dur := time.Since(t0)

	host := r.Host
	url := r.URL.String()
	method := r.Method
	remote := r.RemoteAddr

	revel.INFO.Printf("%v, %s %s%s, %s\n",
		dur, method, host, url, remote)
}

func setJSONMimeType(c *revel.Controller) {
	c.Response.ContentType = "application/json"
}

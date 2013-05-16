// Copyright (c) 2013 The RG Workshop. All rights reserved.
// Use of this source code is governed by a BSD-style license.
// Please feel free to use it.
// More info about the author: www.zeno-l.com

package core

import (
	"appengine"
	"net/http"
)

// Application represents a web app.
type Application interface {
	// Get the app dict.
	GetAppCube() map[string]interface{}
	// Get app default filters.
	GetDefaultFilters() []Filter
	// Set app default filters.
	SetDefaultFilters([]Filter)
	// Get the waitress-handler map.
	GetHandlerMap() map[Waitress]http.Handler
	// Set the waitress-handler map.
	SetHandlerMap(map[Waitress]http.Handler)

	// Get the request-pageCube map.
	GetRequestPageCubeMap() map[*http.Request]map[string]interface{}
	// Get the pageCube directly by the request.
	GetPageCube(*http.Request) map[string]interface{}

	// Forward to another handler by waitress.
	ForwardByWaitress(Waitress, http.ResponseWriter, *http.Request)
}

// Waitress represents every single page.
type Waitress interface {
	// "GET" method job.
	Get(http.ResponseWriter, *http.Request, appengine.Context)
	// "POST" method job.
	Post(http.ResponseWriter, *http.Request, appengine.Context)
	// Get the waitress' filters.
	GetFilters() []Filter
	// Set the waitress' filters.
	SetFilters([]Filter)
	// Get the app.
	GetApp() Application
	// Set the app.
	SetApp(Application)
}

// Filter do the job before the handler's serve function.
type Filter interface {
	// Filte is the function that do the filting job.
	Filte(http.ResponseWriter, *http.Request, appengine.Context) bool
	// Get the app.
	GetApp() Application
	// Set the app that this filter belongs.
	SetApp(Application)
}

// GAEHandler implements the http.Handler interface to serve a web page.
type GAEHandler struct {
	App Application
	Wts Waitress
}

// GAEHandler constuctor.
func NewGAEHandler(waitress Waitress, app Application) *GAEHandler {
	// * init
	self := &GAEHandler{}
	self.App = app

	self.Wts = waitress
	self.Wts.SetApp(app)

	return self
}

func (self *GAEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	c.Debugf("serve a page......")

	// * Get, or build the pageCube.
	var pageCube map[string]interface{}
	var ok bool
	if pageCube, ok = self.App.GetRequestPageCubeMap()[r]; !ok {
		// * Create the pageCube and set to the dict.
		c.Debugf("create page cube......")
		pageCube = make(map[string]interface{})
		self.App.GetRequestPageCubeMap()[r] = pageCube
	}

	// * Do filter stuff, check if it's all valid.
	filters := make([]Filter, 0)
	// Default filters first.
	filters = append(filters, self.App.GetDefaultFilters()...)
	filters = append(filters, self.Wts.GetFilters()...)
	if !isAllValid(filters, w, r, c) {
	} else {
		if r.Method == "GET" {
			self.Wts.Get(w, r, c)
		} else {
			self.Wts.Post(w, r, c)
		}
	}

	// Release page cube.
	c.Debugf("release page cube......")
	delete(self.App.GetRequestPageCubeMap(), r)
}

// This function helps to check all the filters.
func isAllValid(filters []Filter, w http.ResponseWriter, r *http.Request, c appengine.Context) (valid bool) {
	for _, filter := range filters {
		if !filter.Filte(w, r, c) {
			return false
		}
	}
	return true
}

// This function helps to get a build cube.
// Build cube is for the template to build a page.
// Build cube includes the app cube and the page cube.
func GetBuildCube(app Application, r *http.Request) map[string]map[string]interface{} {
	buildCube := make(map[string]map[string]interface{})
	buildCube["AppCube"] = app.GetAppCube()
	buildCube["PageCube"] = app.GetPageCube(r)
	return buildCube
}

// Copyright (c) 2013 The RG Workshop. All rights reserved.
// Use of this source code is governed by a BSD-style license.
// Please feel free to use it.
// More info about the author: www.zeno-l.com

package core

import (
	"net/http"
)

// A simple application.
// It implements the Application interface.
type SimpleApplication struct {
	// The map to store app info.
	AppCube map[string]interface{}
	// The waitress-handler map.
	HandlerMap map[Waitress]http.Handler
	// The request-pageCube map.
	RequestPageCubeMap map[*http.Request]map[string]interface{}
	// App's default filters.
	DefaultFilters []Filter
}

// Constructor of SimpleApplication.
func NewSimpleApplication() *SimpleApplication {
	// * Init.
	self := &SimpleApplication{}
	self.AppCube = make(map[string]interface{})
	self.DefaultFilters = make([]Filter, 0)
	self.HandlerMap = make(map[Waitress]http.Handler)
	self.RequestPageCubeMap = make(map[*http.Request]map[string]interface{})

	return self
}

// Start serve, need the address-waitress map.
func (self *SimpleApplication) StartApp(wtsMap map[string]Waitress) {
	for key, wts := range wtsMap {
		// build handler, and add to handlerMap
		page := NewGAEHandler(wts, self)
		self.GetHandlerMap()[wts] = page
		// serve
		http.Handle(key, page)
	}
}

func (self *SimpleApplication) GetAppCube() map[string]interface{} {
	return self.AppCube
}

func (self *SimpleApplication) SetAppCube(appCube map[string]interface{}) {
	self.AppCube = appCube
}

func (self *SimpleApplication) GetDefaultFilters() []Filter {
	return self.DefaultFilters
}

func (self *SimpleApplication) SetDefaultFilters(filters []Filter) {
	self.DefaultFilters = filters
}

func (self *SimpleApplication) GetHandlerMap() map[Waitress]http.Handler {
	return self.HandlerMap
}

func (self *SimpleApplication) SetHandlerMap(handlerMap map[Waitress]http.Handler) {
	self.HandlerMap = handlerMap
}

func (self *SimpleApplication) GetRequestPageCubeMap() map[*http.Request]map[string]interface{} {
	return self.RequestPageCubeMap
}

func (self *SimpleApplication) GetPageCube(r *http.Request) map[string]interface{} {
	return self.GetRequestPageCubeMap()[r]
}

// Forward to another handler
func (self *SimpleApplication) ForwardByWaitress(wts Waitress, w http.ResponseWriter, r *http.Request) {
	// Find the handler by wts, if not exist, build a new one.
	if _, ok := self.HandlerMap[wts]; !ok {
		self.HandlerMap[wts] = NewGAEHandler(wts, self)
	}
	// Serve.
	self.HandlerMap[wts].ServeHTTP(w, r)
}

// Convenient function to add some default filters to app
func (self *SimpleApplication) AddDefaultFilters(filters ...Filter) {
	self.DefaultFilters = append(self.DefaultFilters, filters...)
}

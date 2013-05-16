// Copyright (c) 2013 The RG Workshop. All rights reserved.
// Use of this source code is governed by a BSD-style license.
// Please feel free to use it.
// More info about the author: www.zeno-l.com

package core

import (
	"appengine"
	"net/http"
)

// A simple waitress.
// It implements the Waitress interface.
type SimpleWaitress struct {
	App     Application
	Filters []Filter
}

// Constructor.
func NewSimpleWaitress() *SimpleWaitress {
	// * Init.
	self := &SimpleWaitress{}
	self.Filters = make([]Filter, 0)

	return self
}

func (self *SimpleWaitress) Get(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	// "GET" method situation.
}

func (self *SimpleWaitress) Post(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	// "POST" method situation.
}

func (self *SimpleWaitress) GetFilters() []Filter {
	return self.Filters
}

func (self *SimpleWaitress) SetFilters(filters []Filter) {
	self.Filters = filters
}

func (self *SimpleWaitress) GetApp() Application {
	return self.App
}

func (self *SimpleWaitress) SetApp(app Application) {
	self.App = app
}

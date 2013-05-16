// Copyright (c) 2013 The RG Workshop. All rights reserved.
// Use of this source code is governed by a BSD-style license.
// Please feel free to use it.
// More info about the author: www.zeno-l.com

package core

import (
	"appengine"
	"net/http"
)

// A simple filter.
// It implements the Filter interface.
type SimpleFilter struct {
	App Application
}

// Constructor.
func NewSimpleFilter(app Application) *SimpleFilter {
	// * Init.
	self := &SimpleFilter{}
	self.App = app

	return self
}

func (self *SimpleFilter) Filte(w http.ResponseWriter, r *http.Request, c appengine.Context) bool {
	return true
}

func (self *SimpleFilter) GetApp() Application {
	return self.App
}

func (self *SimpleFilter) SetApp(app Application) {
	self.App = app
}

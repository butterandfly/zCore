package core

import (
	"appengine"
	"net/http"
)

type Application interface {
	GetAppCube() map[string]interface{}
	SetAppCube(map[string]interface{})
	GetDefaultFilters() []Filter
	SetDefaultFilters([]Filter)
	GetDefaultChefs() []Chef
	SetDefaultChefs([]Chef)
}

type SimpleApplication struct {
	AppCube map[string]interface{}
	// FilterMap      map[string]interface{}
	DefaultFilters []Filter
	DefaultChefs   []Chef
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

func (self *SimpleApplication) GetDefaultChefs() []Chef {
	return self.DefaultChefs
}

func (self *SimpleApplication) SetDefaultChefs(chefs []Chef) {
	self.DefaultChefs = chefs
}

// Waitress interface to perform page
type Waitress interface {
	Get(http.ResponseWriter, *http.Request, appengine.Context)
	Post(http.ResponseWriter, *http.Request, appengine.Context)
	GetFilters() []Filter
	GetChefs() []Chef
	SetAppCube(map[string]interface{})
	GetAppCube() map[string]interface{}
	GetPageCube() map[string]interface{}
	SetPageCube(map[string]interface{})
}

type SimpleWaitress struct {
	AppCube  map[string]interface{}
	PageCube map[string]interface{}
}

func (self *SimpleWaitress) Get(w http.ResponseWriter, r *http.Request, c appengine.Context) {
}

func (self *SimpleWaitress) Post(w http.ResponseWriter, r *http.Request, c appengine.Context) {
}

func (self *SimpleWaitress) GetFilters() []Filter {
	return nil
}

func (self *SimpleWaitress) GetChefs() []Chef {
	return nil
}

func (self *SimpleWaitress) SetAppCube(cube map[string]interface{}) {
	self.AppCube = cube
}

func (self *SimpleWaitress) GetAppCube() map[string]interface{} {
	return self.AppCube
}

func (self *SimpleWaitress) GetPageCube() map[string]interface{} {
	return self.PageCube
}

func (self *SimpleWaitress) SetPageCube(cube map[string]interface{}) {
	self.PageCube = cube
}

type Chef interface {
	Prepare(http.ResponseWriter, *http.Request, appengine.Context)
	SetApp(Application)
}

type SimpleChef struct {
	App Application
}

func (self *SimpleChef) Prepare(w http.ResponseWriter, r *http.Request, c appengine.Context) {

}

func (self *SimpleChef) SetApp(app Application) {
	self.App = app
}

// Filter 
type Filter interface {
	Filte(http.ResponseWriter, *http.Request, appengine.Context) bool
	SetInfoCube(map[string]interface{})
}

type GAEHandler struct {
	App Application
	Wts Waitress
}

func NewGAEHandler(waitress Waitress, app Application) *GAEHandler {
	// * init
	self := &GAEHandler{}
	self.Wts = waitress
	self.App = app

	self.Wts.SetAppCube(self.App.GetAppCube())
	pageCube := make(map[string]interface{})
	pageCube["ActiveNav"] = "home"
	self.Wts.SetPageCube(pageCube)

	return self
}

// this function handles http
func (self *GAEHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// do chef stuff
	chefs := make([]Chef, 0)
	chefs = append(chefs, self.App.GetDefaultChefs()...)
	chefs = append(chefs, self.Wts.GetChefs()...)
	for _, chef := range chefs {
		chef.Prepare(w, r, c)
	}

	// do filter stuff, check if is valid
	filters := make([]Filter, 0)
	filters = append(filters, self.App.GetDefaultFilters()...)
	filters = append(filters, self.Wts.GetFilters()...)
	if !isAllValid(filters, w, r, c) {
		return
	} else {
		if r.Method == "GET" {
			self.Wts.Get(w, r, c)
		} else {
			self.Wts.Post(w, r, c)
		}
	}
}

// this function helps to check all the filters
func isAllValid(filters []Filter, w http.ResponseWriter, r *http.Request, c appengine.Context) (valid bool) {
	for _, filter := range filters {
		if !filter.Filte(w, r, c) {
			return false
		}
	}
	return true
}

// Copyright (c) 2013 The RG Workshop. All rights reserved.
// Use of this source code is governed by a BSD-style license.
// Please feel free to use it.
// More info about the author: www.zeno-l.com

package core

import (
	"appengine"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"
)

// Templer helps to build the page.
type Templer struct {
	// TplMap is a alias-file map.
	// It's helps to get the template file.
	TplMap map[string]interface{}
	// FuncMap is the functions that can be used in
	// the tpl file.
	FuncMap template.FuncMap
}

// Constructor of the Templer.
// "tplMap" is a alias-file map.
// If you don't need costom function in template file, the funcMap can be nil.
func NewTemplerWithMap(funcMap template.FuncMap, tplMap map[string]interface{}) *Templer {
	self := &Templer{}
	self.FuncMap = funcMap
	self.TplMap = tplMap

	return self
}

// Constructor of the Templer.
// "tplFile" is a json file that store alias-file map.
// If you don't need costom function in template file, the funcMap can be nil.
func NewTemplerWithJsonFile(funcMap template.FuncMap, tplFile string) *Templer {
	self := &Templer{}

	b, _ := ioutil.ReadFile(tplFile)
	var v interface{}
	json.Unmarshal(b, &v)
	self.TplMap = v.(map[string]interface{})

	self.FuncMap = funcMap

	return self
}

// This function build the whole page.
func (self *Templer) BuildPage(c appengine.Context, w http.ResponseWriter, infoCube interface{}, tplFiles ...string) (err error) {
	// * Get all template file.
	files := make([]string, 0, len(tplFiles))
	for _, fName := range tplFiles {
		if tplFile, ok := self.TplMap[fName]; ok {
			files = append(files, tplFile.(string))
		}
	}

	// * Build.
	t := template.New("page")
	t = t.Funcs(self.FuncMap)
	t, err = t.ParseFiles(files...)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "page", infoCube)

	return err
}

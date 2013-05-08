package core

import (
	"appengine"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"text/template"
)

type Templer struct {
	TplMap  map[string]interface{}
	FuncMap template.FuncMap
}

func NewTemplerWithJsonFile(funcMap template.FuncMap, tplFile string) *Templer {
	self := &Templer{}

	b, _ := ioutil.ReadFile(tplFile)
	var v interface{}
	json.Unmarshal(b, &v)
	self.TplMap = v.(map[string]interface{})

	self.FuncMap = funcMap

	return self
}

func (self *Templer) BuildPage(c appengine.Context, w http.ResponseWriter, infoCube interface{}, tplFiles ...string) (err error) {
	// * get all tpl file
	files := make([]string, 0, len(tplFiles))
	for _, fName := range tplFiles {
		if tplFile, ok := self.TplMap[fName]; ok {
			files = append(files, tplFile.(string))
		}
	}

	// c.Warningf("%v", files)

	// build
	t := template.New("page")
	t = t.Funcs(self.FuncMap)
	t, err = t.ParseFiles(files...)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "page", infoCube)

	return err
}

package app

import (
	"app/core"
	"appengine"
	"fmt"
	"net/http"
)

// * First, define a page.
var sharedAWaitress = newAWaitress()

type AWaitress struct {
	*core.SimpleWaitress
}

func newAWaitress() *AWaitress {
	self := &AWaitress{core.NewSimpleWaitress()}

	return self
}

func (self *AWaitress) Get(w http.ResponseWriter, r *http.Request, c appengine.Context) {
	fmt.Fprintln(w, "Hello World......")
}

// * Second, set up the address-waitress map 
var wtsMap = map[string]core.Waitress{
	"/": sharedAWaitress,
}

// * Third, start in init function
func init() {
	newApp := core.NewSimpleApplication()

	// Now, enjoy your web!
	newApp.StartApp(wtsMap)
}

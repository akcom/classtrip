package activityfeed

import (
	"appengine"
	"fmt"
	"github.com/mjibson/appstats"
	//"html/template"
	"net/http"
)

func init() {
	http.Handle("/", appstats.NewHandler(Main))
}

func Main(c appengine.Context, w http.ResponseWriter, r *http.Request) {
	// do stuff with c: datastore.Get(c, key, entity)
	fmt.Fprintf(w, "Hello world!")
}

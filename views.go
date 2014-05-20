package classtrip

import (
	"net/http"
	"html/template"
	"sync"
)

var (
	templateNames = []string {"accesserror", "index","settings"}
	templates map[string] *template.Template
	once sync.Once
)

func initTemplates() {
	templates = make(map[string]*template.Template)
	for _, k := range(templateNames) {
		t := template.New("*")
		t = template.Must(t.ParseFiles("templates/master.html"))
		t = template.Must(t.ParseFiles("templates/"+k+".html"))
		templates[k] = t
	}
}

//RenderPage executes the template specified by tmplName and writes the output to w
//data is whatever structures are required by the template
//u is the UserEx structure for the request
func RenderPage(c *Context, w http.ResponseWriter, tmplName string, data interface{}) {
	//initialize the templates
	once.Do(initTemplates)
	//contruct the data to push to the master template
	x := struct {
		Data interface{}
		User  User
	}{ data, c.u }
	
	b := struct {
		Data interface{}
		LogoutURL string
	}{ x, c.logoutURL }
	
	//accesserror does not use the master template, so it is rendered separately
	if tmplName != "accesserror" { 
		templates["index"].ExecuteTemplate(w, "master", b)
	} else {
		templates[tmplName].ExecuteTemplate(w, "accesserror", b)
	}
}
package classtrip

import (
	"net/http"
	"html/template"
	"sync"
	"appengine/user"
	"strings"
)

var (
	templateNames = []string {"accesserror", "index","settings", "calendar"}
	templates map[string] *template.Template
	funcTable template.FuncMap = template.FuncMap { "FirstName" : FirstName, "IsAdmin" : IsAdmin }
	once sync.Once
)

func IsAdmin(ctx *Context) bool {
	return user.IsAdmin(ctx.c)
}

func FirstName(str string) string {
	return strings.Split(str, " ")[0]
}

func initTemplates() {
	templates = make(map[string]*template.Template)
	for _, k := range(templateNames) {
		t := template.New("*").Funcs(funcTable)
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
		Ctx *Context
	}{ data, c.u, c }
	
	b := struct {
		Data interface{}
		Active string
		LogoutURL string
	}{ x, tmplName, c.logoutURL }
	
	//accesserror does not use the master template, so it is rendered separately
	if tmplName != "accesserror" { 
		templates[tmplName].ExecuteTemplate(w, "master", b)
	} else {
		templates[tmplName].ExecuteTemplate(w, "accesserror", b)
	}
}
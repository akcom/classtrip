package classtrip

import (
	"appengine"
	"appengine/user"
	"appengine/datastore"
	"github.com/mjibson/appstats"
	"net/http"
	"html/template"
	"strings"
	"time"
	"strconv"
)

type User struct {
	Email string
	Name string
	CreateTime time.Time
}

type Context struct {
	u User
	c appengine.Context
	cu *user.User
	logoutURL string
}

func init() {
	http.Handle("/", appstats.NewHandler(Secure(Main)))
	http.Handle("/main", appstats.NewHandler(Secure(Main)))
	http.Handle("/settings", appstats.NewHandler(Secure(Settings)))
	http.Handle("/admin", appstats.NewHandler(Secure(AdminMain)))
	http.Handle("/calendar", appstats.NewHandler(Secure(Calendar)))
	http.Handle("/calendar/", appstats.NewHandler(Secure(Calendar)))
}

//Secure wraps a Handler functions and includes functionality to make sure
//the user is from university of maryland and also pulls their User structure
type StatsHandler func (appengine.Context, http.ResponseWriter, *http.Request)
type CustomHandler func (*Context, http.ResponseWriter, *http.Request)

func Secure(handler CustomHandler) StatsHandler {
	return func(c appengine.Context, w http.ResponseWriter, r *http.Request) {
		//store a context which allows us to avoid repeated calls to user.Current()
		//user.LogoutURL() etc
		ctx := new(Context)
		ctx.c = c
		ctx.cu = user.Current(c)
		ctx.logoutURL, _ = user.LogoutURL(c, "/main")
		
		//make sure the user is a university of maryland student before continuing
		if !strings.HasSuffix(ctx.cu.Email, "@umaryland.edu") {
			//return status unauthorized
			w.WriteHeader(http.StatusUnauthorized)
			//the access denied template gives the user the option to logout and try again
			RenderPage(ctx, w, "accesserror", ctx.cu)
			return
		}
		
		//the user is from UMB, so now try to pull their user record
		//the user record is created if it does not exist
		
		//the root key is the key under which all User{}'s live
		userRK := userRootKey(c)

		//create the key based on the username
		userKey := datastore.NewKey(ctx.c, "User", ctx.cu.Email, 0, userRK)
		
		//check if there is already an entry for the user, if not, create one.
		if err := datastore.Get(ctx.c, userKey, &ctx.u); err != nil {
			ctx.u.Email = ctx.cu.Email
			ctx.u.CreateTime = time.Now()
			ctx.u.Name = ""
		
			if _,err := datastore.Put(c, userKey, &ctx.u); err != nil {
				c.Errorf("Unable to put user key: %s", err.Error())
			}
		}
		handler(ctx, w, r)
	}
}

//userKey returns the root User key under which all User{} structures lay
func userRootKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "User", "default_users", 0, nil)
}

//userKey returns the key for a given email
func userKey(c appengine.Context, email string) *datastore.Key {
	c.Debugf("Creating key for %d", email)
	return datastore.NewKey(c, "User", email, 0, userRootKey(c))
}

//success returns the HTML "text-success"  wrapped message
func success(message string) template.HTML {
	return template.HTML(`<p class="bg-success">` + message + `</p>`)
}

//failure returns the HTML "text-danger"  wrapped message
func failure(message string) template.HTML {
	return template.HTML(`<p class="bg-danger">` + message + `</p>`)
}

func Main(c *Context, w http.ResponseWriter, r *http.Request) {
	RenderPage(c, w, "index", nil)
}

func Settings(c *Context, w http.ResponseWriter, r *http.Request) {
	//helper function
	render := func (str template.HTML) {
		RenderPage(c, w, "settings", str)
	}
	if (r.Method == "GET") {
		render("")
	} else {
		r.ParseForm();
		param,ok := r.Form["fullName"]
		if !ok {
			RenderPage(c, w, "settings", failure("Error updating settings: fullname not submitted"))
			return
		}
		name := param[0]
		//check to make sure the name is valid (at least two words)
		if !strings.Contains(name, " ") {
			render(failure("Make sure you entered your full name (first and last)"))
			return
		}
		//get the user's datastore key
		k := userKey(c.c, c.u.Email)
		c.u.Name = name
		datastore.Put(c.c, k, &c.u)
		RenderPage(c, w, "settings", success("Settings updated"))
	}
}


func Calendar(c *Context, w http.ResponseWriter, r *http.Request) {
	var (
		cal string
		url string = r.URL.Path
		t time.Time = time.Now()
		month time.Month
		year int
	)
	//parse the URL to extract the month and the year
	
	//check to see if a month and year was provided and is in the proper format
	//URL should be in format of calendar/month/year
	//where month is a 1-2 digit integer and year is a four digit integer
	idx := strings.LastIndex(url, "calendar/")
	
	//default is the current month and year in case the the request is improperly formatted
	month = t.Month()
	year = t.Year()
	if idx != -1 { 
		//now we have to determine if the URL is properly formatted
		//ex: calendar/01/2014 for Jan 2014
		
		//split them by the '/'
		strs := strings.Split(url[idx+9:], "/")
		//check to make sure we have two items (month and a year) and they are the proper length
		if len(strs) == 2 && len(strs[1]) == 4 && (len(strs[0]) == 1 || len(strs[0]) == 2) {
			//format is valid
			if yearInt,err := strconv.Atoi(strs[1]); err == nil && yearInt >= year {
				if monthInt,err := strconv.Atoi(strs[0]); err == nil {
					//cant go back in time
					if yearInt > year || monthInt > int(month) {
						month = time.Month(monthInt)
						year = yearInt
					}
				}
			}
		}
	}
	cal = GenCalendar(month, year)
	RenderPage(c, w, "calendar", template.HTML(cal))
}

func AdminMain(c *Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world"))
}
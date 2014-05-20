package classtrip

import (
	"appengine"
	"appengine/user"
	"appengine/datastore"
	"fmt"
	"github.com/mjibson/appstats"
	"net/http"
	"strings"
	"time"
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
	http.Handle("/main", appstats.NewHandler(Secure(Main)))
	http.Handle("/settings/", appstats.NewHandler(Secure(Settings)))
	http.Handle("/admin/", appstats.NewHandler(Secure(AdminMain)))	
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
		ctx.u.Email = ctx.cu.Email
		
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
		if err := datastore.Get(ctx.c, userKey, ctx.u); err == nil {
			panic("Unable to get user key")
		}
		
		//is the entry does not exist...
		if ctx.u.Email == "" { 
			ctx.u.Email = ctx.cu.Email
			ctx.u.CreateTime = time.Now()
			ctx.u.Name = "Unknown"
		
			datastore.Put(c, userKey, ctx.u)
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
	return datastore.NewKey(c, "User", email, 0, userRootKey(c))
}

func Main(c *Context, w http.ResponseWriter, r *http.Request) {
	RenderPage(c, w, "index", nil)
}

func Settings(c *Context, w http.ResponseWriter, r *http.Request) {
	RenderPage(c, w, "settings", nil)
}

func AdminMain(c *Context, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello admin!")
}
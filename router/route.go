package router

import (
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
	"github.com/nayan9800/shorturl/data"
)

var store = sessions.NewCookieStore([]byte("SECRET"))

func Method(h http.HandlerFunc, method string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(rw, "Bad Request", http.StatusBadRequest)
			return
		}
		h(rw, r)
	}
}
func isAuthenticated(r *http.Request) bool {
	sess, _ := store.Get(r, "shorturl-auth")
	auth, ok := sess.Values["authenticated"].(bool)
	return !ok || !auth
}

func Index(rw http.ResponseWriter, r *http.Request) {

	code := r.URL.RawQuery
	redir := data.GetUrl(code)
	if redir.Url == "" {
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	}
	http.Redirect(rw, r, redir.Url, http.StatusFound)
}
func Dashboard(rw http.ResponseWriter, r *http.Request) {

	sess, _ := store.Get(r, "shorturl-auth")
	if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
		http.Error(rw, "Forbiddan", http.StatusForbidden)
		return
	}
	id := sess.Values["Id"].(int)
	User := data.GetUserById(id)
	dash, err := template.ParseFiles("template/dashboard.html")
	if err != nil {
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}
	dash.Execute(rw, User)
}
func AddService(rw http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Error(rw, "Forbiden", http.StatusForbidden)
		return
	}
	sess, _ := store.Get(r, "shorturl-auth")
	id, _ := sess.Values["Id"].(int)
	switch {
	case r.Method == "GET":
		l, err := template.ParseFiles("template/addservice.html")
		if err != nil {
			log.Println(err.Error())
		}
		l.Execute(rw, nil)

	case r.Method == "POST":
		r.ParseForm()
		newurl := r.FormValue("new-url")
		data.InsertSevice(newurl, id)
		http.Redirect(rw, r, "/dashboard", http.StatusFound)

	default:
		http.Error(rw, "Bad request", http.StatusBadRequest)

	}
}
func DeleteService(rw http.ResponseWriter, r *http.Request) {
	if isAuthenticated(r) {
		http.Error(rw, "Forbiden", http.StatusForbidden)
		return
	}
	sess, _ := store.Get(r, "shorturl-auth")
	id, _ := sess.Values["Id"].(int)
	switch {
	case r.Method == "GET":
		l, err := template.ParseFiles("template/delete.html")
		if err != nil {
			log.Println(err.Error())
		}
		u := data.GetUserById(id)
		l.Execute(rw, u)

	case r.Method == "POST":
		r.ParseForm()
		scode := r.FormValue("service")
		data.DeleteServiceByCode(scode)
		http.Redirect(rw, r, "/dashboard", http.StatusFound)

	default:
		http.Error(rw, "Bad request", http.StatusBadRequest)
	}
}
func Logout(rw http.ResponseWriter, r *http.Request) {
	sess, _ := store.Get(r, "shorturl-auth")
	sess.Options.MaxAge = -1
	sess.Save(r, rw)
	http.Redirect(rw, r, "/login", http.StatusFound)

}
func Loginpage(rw http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == "GET":
		l, err := template.ParseFiles("template/login.html")
		if err != nil {
			log.Println(err.Error())
		}
		l.Execute(rw, nil)

	case r.Method == "POST":
		r.ParseForm()
		email := r.FormValue("Emailid")
		password := r.FormValue("Password")
		id, err := data.AuthUser(email, password)
		if err == data.DataPasswordNotMatch {
			http.Error(rw, "Invalid credentials", http.StatusUnauthorized)
			return
		} else if err != nil {
			http.Error(rw, "Not Found", http.StatusNotFound)
			return
		}
		sess, _ := store.Get(r, "shorturl-auth")
		sess.Values["authenticated"] = true
		sess.Values["Id"] = id
		sess.Save(r, rw)
		http.Redirect(rw, r, "/dashboard", http.StatusFound)

	default:
		http.Error(rw, "BadRequest", http.StatusBadRequest)

	}
}
func SignupPage(rw http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == "GET":
		l, err := template.ParseFiles("template/signup.html")
		if err != nil {
			log.Println(err.Error())
		}
		l.Execute(rw, nil)

	case r.Method == "POST":
		r.ParseForm()
		email := r.FormValue("Emailid")
		password := r.FormValue("Password")
		name := r.FormValue("name")
		u := data.User{Name: name, Password: password, Email: email}
		err := u.InsertUser()
		if err != nil {
			http.Error(rw, "Try login", http.StatusNotFound)
		}
		nu, _ := data.GetUser(u.Email)
		sess, _ := store.Get(r, "shorturl-auth")
		sess.Values["authenticated"] = true
		sess.Values["Id"] = nu.ID
		sess.Save(r, rw)
		http.Redirect(rw, r, "/dashboard", http.StatusFound)

	default:
		http.Error(rw, "BadRequest", http.StatusBadRequest)

	}
}

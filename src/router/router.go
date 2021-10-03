package router

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"rn.com/src/config"
)

type IndexHandler struct{}

func (h IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "./view/index.html")
	}
}

type AdminHandler struct{}

func (h AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if err := r.ParseForm(); err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		pass := r.FormValue("rn-id")

		hash := sha256.New()
		hash.Write([]byte(pass))
		md := hash.Sum(nil)
		pass = hex.EncodeToString(md)
		if pass == config.PASS {
			cookie := http.Cookie{Name: "RN_TOKEN", Value: config.TOKEN, HttpOnly: true}
			http.SetCookie(w, &cookie)
			http.Redirect(w, r, "/task", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

type TaskHandler struct{}

func (h TaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		token, err := r.Cookie("RN_TOKEN")
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		if token.Value == config.TOKEN {
			http.ServeFile(w, r, "./view/task.html")
			return
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}

func Get() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", &IndexHandler{})
	mux.Handle("/check-admin", &AdminHandler{})
	mux.Handle("/task", &TaskHandler{})
	return mux
}

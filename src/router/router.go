package router

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/png"
	"net/http"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
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
	token, err := r.Cookie("RN_TOKEN")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if token.Value == config.TOKEN {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, "./view/task.html")
		}
		return
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

type BaramHandler struct{}
type RnPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (h BaramHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := r.Cookie("RN_TOKEN")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if token.Value == config.TOKEN {
		switch r.Method {
		case "GET":
			bounds := image.Rectangle{Min: image.Point{0, 28}, Max: image.Point{788, 640}}
			img, _ := screenshot.CaptureRect(bounds)
			// if err != nil {
			// 	panic(err)
			// }
			w.Header().Set("Content-Type", "image/png")
			png.Encode(w, img)
		case "POST":
			p := new(RnPoint)
			json.NewDecoder(r.Body).Decode(p)
			// fmt.Println(p.X, p.Y+28)
			robotgo.MoveMouse(p.X, p.Y+28)
			time.Sleep(50 * time.Millisecond)
			robotgo.MouseClick("left", true)
			time.Sleep(100 * time.Millisecond)
		}
		return
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func Get() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", &IndexHandler{})
	mux.Handle("/check-admin", &AdminHandler{})
	mux.Handle("/task", &TaskHandler{})
	mux.Handle("/baram", &BaramHandler{})
	return mux
}

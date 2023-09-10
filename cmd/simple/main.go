package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"text/template"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type PageShort struct {
	URL     string
	Success bool
	Error   error
	Short   string
}

type Page404 struct {
	Item string
}

type Page500 struct {
	Message string
}

func main() {
	tmpl := template.Must(template.ParseGlob("cmd/simple/templates/*.html"))
	shortner := UrlShortner{
		urls: make(map[string]string),
	}
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/s/{hash}", func(w http.ResponseWriter, r *http.Request) {
		hash := chi.URLParam(r, "hash")
		log.Println(hash)
		url, ok := shortner.Get(hash)
		if !ok {
			w.WriteHeader(404)
			tmpl.ExecuteTemplate(w, "404.html", Page404{Item: "Short"})
			return
		}
		http.Redirect(w, r, url, http.StatusFound)
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "cmd/simple/index.html")
	})

	r.Post("/new", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		url := r.FormValue("url")
		hash := shortner.Generate(url)
		short := fmt.Sprintf("http://%s/s/%s", r.Host, hash)
		data := PageShort{
			URL:     url,
			Success: true,
			Short:   short,
			Error:   nil,
		}
		var tpl bytes.Buffer
		err := tmpl.ExecuteTemplate(&tpl, "short.html", data)
		if err != nil {
			log.Println(err)
			data.Error = err
			w.WriteHeader(500)
			tmpl.ExecuteTemplate(w, "500.html", Page500{Message: "Could not render template"})
			return
		}

		w.Write(tpl.Bytes())

	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		tmpl.ExecuteTemplate(w, "404.html", Page404{Item: "Page"})
	})

	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		fs := http.StripPrefix("/static", http.FileServer(http.Dir("dist")))
		fs.ServeHTTP(w, r)
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}

	http.ListenAndServe(":3000", r)
}

type UrlShortner struct {
	//[short]url
	urls map[string]string
}

func (shortner *UrlShortner) Generate(url string) (short string) {
	hash := shortner.short(charset)
	shortner.urls[hash] = url
	return hash
}

func (shortner *UrlShortner) Get(short string) (url string, found bool) {
	url, found = shortner.urls[short]

	return url, found
}

func (shortner *UrlShortner) HasShort(short string) bool {
	_, ok := shortner.urls[short]
	return ok
}

func (shortner *UrlShortner) short(charset string) string {
	var hash string
	for {
		b := make([]byte, 6)
		for i := range b {
			b[i] = charset[rand.Intn(len(charset))]
		}
		newHash := string(b)
		if ok := shortner.HasShort(hash); !ok {
			log.Println("New hash:", newHash)
			hash = newHash
			break
		}
	}

	return hash
}

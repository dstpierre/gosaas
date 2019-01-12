package gosaas

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

var (
	pageTemplates *template.Template
	languagePacks map[string]map[string]string
)

func init() {
	loadTemplates()
	loadLanguagePacks()
}

func loadTemplates() {
	var tmpl []string

	files, err := ioutil.ReadDir("./templates")
	if err != nil {
		if os.IsNotExist(err) == false {
			log.Fatal("unable to load templates", err)
		}
		return
	}

	for _, f := range files {
		tmpl = append(tmpl, path.Join("./templates", f.Name()))
	}

	t, err := template.New("").Funcs(template.FuncMap{
		"translate":  Translate,
		"translatef": Translatef,
		"money": func(amount int) string {
			m := float64(amount) / 100.0
			return fmt.Sprintf("%.2f $", m)
		},
	}).ParseFiles(tmpl...)

	if err != nil {
		log.Fatal("error while parsing templates", err)
	}

	pageTemplates = t
}

// ServePage will render and respond with an HTML template.ServePage
//
// HTML templates should be saved into a directory named templates.ServePage
//
// Example usage:
//
// 	func handler(w http.ResponseWriter, r *http.Request) {
// 		data := HomePage{Title: "Hello world!"}
// 		gosaas.ServePage(w, r, "index.html", data)
// 	}
func ServePage(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	t := pageTemplates.Lookup(name)

	if err := t.Execute(w, data); err != nil {
		fmt.Println("error while rendering the template ", err)
	}

	logRequest(r, http.StatusOK)
}

func loadLanguagePacks() {
	languagePacks = make(map[string]map[string]string)

	files, err := ioutil.ReadDir("./languagepacks")
	if err != nil {
		log.Println("unable to load language packs: ", err)
		return
	}

	var pack = new(struct {
		Language string `json:"lang"`
		Keys     []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"keys"`
	})

	for _, f := range files {
		b, err := ioutil.ReadFile(path.Join("./languagepacks", f.Name()))
		if err != nil {
			log.Fatal("unable to read language pack: ", f.Name(), ": ", err)
		}

		if err := json.Unmarshal(b, &pack); err != nil {
			log.Fatal("unable to parse language pack: ", f.Name(), ": ", err)
		}

		values := make(map[string]string)
		for _, k := range pack.Keys {
			values[k.Key] = k.Value
		}

		languagePacks[pack.Language] = values
	}
}

// Translate finds a key in a language pack file (saved in directory named languagepack)
// and return the value as template.HTML so it's safe to use HTML inside the language pack file.Translate
//
// The language pack file are simple JSON file named lng.json like en.json:
//
// 	{
// 		"lang": "en",
// 		"keys": [
// 			{"key": "landing-title", "value": "Welcome to my site"}
// 		]
// 	}
func Translate(lng, key string) template.HTML {
	if s, ok := languagePacks[lng][key]; ok {
		return template.HTML(s)
	}
	return template.HTML(fmt.Sprintf("key %s not found", key))
}

// Translatef finds a translation key and substitute the formatting parameters.
func Translatef(lng, key string, a ...interface{}) string {
	if s, ok := languagePacks[lng][key]; ok {
		return fmt.Sprintf(s, a...)
	}
	return fmt.Sprintf("key %s not found", key)
}

// BUG(dom): This needs more thinking...
func ExtractPageAndFilter(r *http.Request) (page int, filter string) {
	p := r.URL.Query().Get("p")
	if len(p) > 0 {
		i, err := strconv.Atoi(p)
		if err == nil {
			page = i
		}
	}

	filter = r.URL.Query().Get("filter")

	return
}

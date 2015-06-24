package ui

import(
	"net/http"
	"html/template"
	"path"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/apenella/verdmell", http.StatusFound)
}

func WebUI(w http.ResponseWriter, r *http.Request){

	layout := path.Join("ui","templates", "layout.html")
	index := path.Join("ui","templates", "index.html")

 	tmpl, err := template.ParseFiles(layout, index)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//#######################################################################################################
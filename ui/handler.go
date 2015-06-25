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

	index := path.Join("ui","templates", "index.html")
	header := path.Join("ui","templates", "header.html")
	content := path.Join("ui","templates", "content.html")
	footer := path.Join("ui","templates", "footer.html")

 	tmpl, err := template.ParseFiles(index,header,content,footer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//#######################################################################################################
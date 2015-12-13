package ui

import(
	"net/http"
	"html/template"
	"path"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/ui", http.StatusFound)
}

func WebUI(w http.ResponseWriter, r *http.Request){

	webui := GetUI()
	cluster := webui.nodes

	index := path.Join("ui","html", "index.html")
	scripts := path.Join("ui","scripts", "scripts.js")
	style := path.Join("ui","style", "verdmell.css")
	header := path.Join("ui","html", "header.html")
	content := path.Join("ui","html", "content.html")
	footer := path.Join("ui","html", "footer.html")


 	tmpl, err := template.ParseFiles(index,scripts,style,header,content,footer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, cluster); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//#######################################################################################################
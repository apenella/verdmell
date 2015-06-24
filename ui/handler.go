package ui

import(
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://github.com/apenella/verdmell", http.StatusFound)
}

//#######################################################################################################
package api

import(
	"verdmell/check"
	"net/http"
)

//
//# Index: is the handler that manages the root api request
func Index(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::Index)")
	http.Redirect(w, r, "/api/run", http.StatusFound)
}
//
//# StartCheckSystem: is the handler that manages the start checks system request
func StartCheckSystem(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::StartCheckSystem)")
	checks := box.GetObject(CHECKS).(*check.CheckSystem)
	//vars := mux.Vars(r)
	//check := vars["check"]

	if err := checks.StartCheckSystem(nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	GetAllServices(w,r)
}

//#######################################################################################################
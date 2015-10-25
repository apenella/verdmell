package api

import(
	"verdmell/check"
	"verdmell/service"
	"net/http"
)

//
//# StartCheckEngine: is the handler that manages the start checks system request
func StartCheckEngine(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::StartCheckEngine)")
	checks := box.GetObject(CHECKS).(*check.CheckEngine)
	//vars := mux.Vars(r)
	//check := vars["check"]

	if err := checks.StartCheckEngine(nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	GetAllServices(w,r)
}

//
//# GetNodeStatus: write the specific service's data to response writer
func GetNodeStatus(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::GetService)")
	services := box.GetObject(SERVICES).(*service.ServiceEngine)
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(services.GetService(env.Context.Service))
}

//#######################################################################################################
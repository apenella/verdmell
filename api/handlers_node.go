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
	checks := env.GetCheckEngine().(*check.CheckEngine)
	//vars := mux.Vars(r)
	//check := vars["check"]

	if err := checks.StartCheckEngine(nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	apiWriter(GetAllServices)
}

//
//# GetNodeStatus: write the specific service's data to response writer
func GetNodeStatus(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::GetService)")
	services := env.GetServiceEngine().(*service.ServiceEngine)

	_,json := services.GetService(env.Context.Service)
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

//#######################################################################################################
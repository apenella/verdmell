package api

import(
	"verdmell/check"
	"verdmell/service"
	"net/http"
)

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

//
//# GetNodeStatus: write the specific service's data to response writer
func GetNodeStatus(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetService)")
	services := box.GetObject(SERVICES).(*service.ServiceSystem)
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(services.GetService(env.Context.Service))
}

//#######################################################################################################
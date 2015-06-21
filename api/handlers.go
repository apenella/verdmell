package api

import(
	"verdmell/check"
	"verdmell/service"
	"net/http"

	"github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::Index)")
	http.Redirect(w, r, "/api/run", http.StatusOK)
}

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


func GetAllChecks(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetAllChecks)")
	checks := box.GetObject(CHECKS).(*check.CheckSystem)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetAllChecks())	
}

func GetCheck(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetCheck)")
	checks := box.GetObject(CHECKS).(*check.CheckSystem)
	vars := mux.Vars(r)
	check := vars["check"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetCheck(check))
}

func GetAllServices(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetAllServices)")
	services := box.GetObject(SERVICES).(*service.ServiceSystem)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(services.GetAllServices())
}

func GetService(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetService)")
	services := box.GetObject(SERVICES).(*service.ServiceSystem)
	vars := mux.Vars(r)
	service := vars["service"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(services.GetService(service))
}

func GetAllSamples(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetAllSamples)")
	checks := box.GetObject(CHECKS).(*check.CheckSystem)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetAllSamples())	
}

func GetSample(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetSample)")
	checks := box.GetObject(CHECKS).(*check.CheckSystem)
	vars := mux.Vars(r)
	sample := vars["sample"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetSampleForCheck(sample))
}
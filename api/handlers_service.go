package api

import(
	"verdmell/service"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllServices: write all services' data to response writer
func GetAllServices(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetAllServices)")
	services := box.GetObject(SERVICES).(*service.ServiceEngine)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(services.GetAllServices())
}
//
//# GetService: write the specific service's data to response writer
func GetService(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetService)")
	services := box.GetObject(SERVICES).(*service.ServiceEngine)
	vars := mux.Vars(r)
	service := vars["service"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(services.GetService(service))
}

//#######################################################################################################
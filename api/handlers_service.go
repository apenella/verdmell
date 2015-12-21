package api

import(
	"verdmell/service"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllServices: write all services' data to response writer
func GetAllServices(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::GetAllServices)")
	services := env.GetServiceEngine().(*service.ServiceEngine)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	_,json := services.GetAllServices()

	w.Write(json)
}
//
//# GetService: write the specific service's data to response writer
func GetService(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::GetService)")
	services := env.GetServiceEngine().(*service.ServiceEngine)
	vars := mux.Vars(r)
	service := vars["service"]

	_,json := services.GetService(service)


	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

//#######################################################################################################
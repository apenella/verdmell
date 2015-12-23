package api

import(
	"verdmell/service"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllServices: write all services' data to response writer
func GetAllServices(r *http.Request) (error, []byte) {
	env.Output.WriteChDebug("(ApiEngine::GetAllServices)")
	services := env.GetServiceEngine().(*service.ServiceEngine)

	return services.GetAllServices()
}
//
//# GetService: write the specific service's data to response writer
func GetService(r *http.Request) (error, []byte) {
	env.Output.WriteChDebug("(ApiEngine::GetService)")
	services := env.GetServiceEngine().(*service.ServiceEngine)
	vars := mux.Vars(r)
	service := vars["service"]

	return services.GetService(service)
}

//#######################################################################################################
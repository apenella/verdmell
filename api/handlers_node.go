package api

import(
	"errors"
	"verdmell/check"
	"verdmell/service"
	"net/http"
)

//
//# GetNodeStatus: write the specific service's data to response writer
func GetNodeInfo(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetNodeInfo)")
	return env.GetNodeInfo()
}
//
//# GetNodeStatus: write the specific service's data to response writer
func GetNodeStatus(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetNodeStatus)")
	services := env.GetServiceEngine().(*service.ServiceEngine)

	return services.GetService(env.Context.Service)
}
//
//# StartCheckEngine: is the handler that manages the start checks system request
func StartCheckEngine(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::StartCheckEngine)")
	checks := env.GetCheckEngine().(*check.CheckEngine)
	
	if err := checks.StartCheckEngine(nil); err != nil {
		return errors.New("(ApiEngine::StartCheckEngine) "+err.Error()),nil
	}
	return GetNodeStatus(r)
}

//#######################################################################################################
package api

import(
	"verdmell/check"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllSamples: write all samples' data to response writer
func GetAllSamples(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetAllSamples)")
	checks := box.GetObject(CHECKS).(*check.CheckEngine)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetAllSamples())	
}
//
//# GetSample: write the specific sample's data to response writer
func GetSample(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetSample)")
	checks := box.GetObject(CHECKS).(*check.CheckEngine)
	vars := mux.Vars(r)
	sample := vars["sample"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetSampleForCheck(sample))
}

//#######################################################################################################
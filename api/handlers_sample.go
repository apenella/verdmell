package api

import(
	"verdmell/sample"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllSamples: write all samples' data to response writer
func GetAllSamples(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::GetAllSamples)")
	checks := env.GetSampleEngine().(*sample.SampleEngine)
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetAllSamples())	
}
//
//# GetSample: write the specific sample's data to response writer
func GetSample(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::GetSample)")
	
	samples := env.GetSampleEngine().(*sample.SampleEngine)
	vars := mux.Vars(r)
	check := vars["sample"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	w.Write(samples.GetSampleForCheck(check))
}

//#######################################################################################################
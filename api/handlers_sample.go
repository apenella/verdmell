package api

import(
	"verdmell/sample"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllSamples: write all samples' data to response writer
func GetAllSamples(r *http.Request) (error, []byte) {
	env.Output.WriteChDebug("(ApiEngine::GetAllSamples)")
	samples := env.GetSampleEngine().(*sample.SampleEngine)
	
	return samples.GetAllSamples()
}
//
//# GetSample: write the specific sample's data to response writer
func GetSample(r *http.Request) (error, []byte) {
	env.Output.WriteChDebug("(ApiEngine::GetSample)")
	
	samples := env.GetSampleEngine().(*sample.SampleEngine)
	vars := mux.Vars(r)
	check := vars["sample"]

	return samples.GetSampleForCheck(check)
}

//#######################################################################################################
package api

import(
	"verdmell/check"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllChecks: write all checks' data to response writer
func GetAllChecks(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::GetAllChecks)")
	checks := env.GetCheckEngine().(*check.CheckEngine)
	//checks := box.GetObject(CHECKS).(*check.CheckEngine)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetAllChecks())	
}
//
//# GetCheck: write the specific check's data to response writer
func GetCheck(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::GetCheck)")
	checks := env.GetCheckEngine().(*check.CheckEngine)
	vars := mux.Vars(r)
	check := vars["check"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetCheck(check))
}

//#######################################################################################################
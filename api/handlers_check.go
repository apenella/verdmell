package api

import(
	"verdmell/check"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllChecks: write all checks' data to response writer
func GetAllChecks(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetAllChecks)")
	checks := box.GetObject(CHECKS).(*check.CheckSystem)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetAllChecks())	
}
//
//# GetCheck: write the specific check's data to response writer
func GetCheck(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiSystem::GetCheck)")
	checks := box.GetObject(CHECKS).(*check.CheckSystem)
	vars := mux.Vars(r)
	check := vars["check"]
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(checks.GetCheck(check))
}

//#######################################################################################################
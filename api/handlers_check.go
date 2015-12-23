package api

import(
	"verdmell/check"
	"net/http"
	"github.com/gorilla/mux"
)

//
//# GetAllChecks: write all checks' data to response writer
func GetAllChecks(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetAllChecks)")
	checks := env.GetCheckEngine().(*check.CheckEngine)

	return checks.GetAllChecks()
}
//
//# GetCheck: write the specific check's data to response writer
func GetCheck(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetCheck)")
	checks := env.GetCheckEngine().(*check.CheckEngine)
	vars := mux.Vars(r)
	check := vars["check"]
	
	return checks.GetCheck(check)
}

//#######################################################################################################
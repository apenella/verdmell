package api

import(
	"net/http"
)

//
//# Index: is the handler that manages the root api request
func Index(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::Index)")
	http.Redirect(w, r, "/api/node", http.StatusFound)
}

//
//# GetCluster: is the handler that return all cluster nodes
// func GetCluster(w http.ResponseWriter, r *http.Request) {
// 	env.Output.WriteChDebug("(ApiEngine::GetClusterNodes)")
	
// 	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(env.GetCluster())
// }

//#######################################################################################################
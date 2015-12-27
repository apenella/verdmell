package api

import(
	"verdmell/cluster"
	"net/http"
)

//
//# GetAllChecks: write all checks' data to response writer
func GetClusterNodes(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetAllChecks)")
	cluster := env.GetClusterEngine().(*cluster.ClusterEngine)

	return cluster.GetClusterNodes()
}

//#######################################################################################################
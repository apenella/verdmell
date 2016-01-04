package api

import(
	"verdmell/cluster"
	"net/http"
)

//
//# GetClusterInfo: write whole cluster's information to response writer
func GetClusterInfo(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetClusterInfo)")
	cluster := env.GetClusterEngine().(*cluster.ClusterEngine)

	return cluster.GetClusterInfo()
}

//
//# GetClusterNodes: write all checks' data to response writer
func GetClusterNodes(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetClusterNodes)")
	cluster := env.GetClusterEngine().(*cluster.ClusterEngine)

	return cluster.GetClusterNodes()
}

//#######################################################################################################
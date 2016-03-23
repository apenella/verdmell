package api

import(
	//"encoding/binary"
	//"encoding/json"
	"io/ioutil"
	"net/http"
	"verdmell/cluster"
)

//
//# GetCluster: write whole cluster's information to response writer
func GetCluster(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetCluster)")
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
//
//# GetClusterServices: write all checks' data to response writer
func GetClusterServices(r *http.Request) (error,[]byte) {
	env.Output.WriteChDebug("(ApiEngine::GetClusterServices)")
	cluster := env.GetClusterEngine().(*cluster.ClusterEngine)

	return cluster.GetClusterServices()
}
//
//# PutClusterStatus: write all checks' data to response writer
func PutCluster(w http.ResponseWriter, r *http.Request) {
	env.Output.WriteChDebug("(ApiEngine::PutClusterStatus)")
	cluster := env.GetClusterEngine().(*cluster.ClusterEngine)
	//
	//curl -X PUT -d 'hola' http://0.0.0.0:5497/api/cluster
	if body, err := ioutil.ReadAll(r.Body); err != nil {
		env.Output.WriteChDebug("(ApiEngine::PutClusterStatus) Could not read the request body")
	} else {
		cluster.ReceiveData(body)
	}
}

//#######################################################################################################
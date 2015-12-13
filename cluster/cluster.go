/*
Cluster Engine management

The package 'cluster' is used by verdmell to manage the cluster. 

=> Is known as a cluster a set of nodes.

-ClusterEngine
-Cluster
-ClusterNode
*/

package cluster

import (
  "errors"
	"verdmell/utils"
)

//#
//#
//# Cluster struct:
//# Cluster is a set of nodes
type Cluster struct{
	//map to store the nodes that belong to the cluster
  Nodes map[string]*ClusterNode `json:"nodes"`
  //map to store the services that belong to the cluster
  Services map[string] string `json:"services"` 
}

//
//# NewCluster: return a CheckEngine instance to be run
func NewCluster() (error, *Cluster) {
// func NewCluster(n map[string]*ClusterNode, s map[string]string) (error, *Cluster) {
	env.Output.WriteChDebug("(Cluster::NewCluster)")
  cluster := new(Cluster)

  // cluster.SetNodes(n)
  // cluster.SetServices(s)

	return nil, cluster
}

//#
//# Getters and Setters
//#----------------------------------------------------------------------------------------

//
//# SetNodes: set attribute from Cluster
func (c *Cluster) SetNodes(nodes map[string]*ClusterNode) {
  env.Output.WriteChDebug("(Cluster::SetNodes) Set Nodes' value")
  env.Output.WriteChDebug(nodes)
  c.Nodes = nodes
}

//
//# GetNodes: get attribute from Cluster
func (c *Cluster) GetNodes() map[string]*ClusterNode {
  env.Output.WriteChDebug("(Cluster::GetNodes) Get Nodes' value")
  env.Output.WriteChDebug(c.Nodes)
  return c.Nodes
}

//
//# SetServices: set attribute from Cluster
func (c *Cluster) SetServices(services map[string] string) {
  env.Output.WriteChDebug("(Cluster::SetServices) Set Services' value")
  c.Services = services
}

//
//# GetServices: get attribute from Cluster
func (c *Cluster) GetServices() map[string] string {
  env.Output.WriteChDebug("(Cluster::GetNodes) Get Services' value")
  return c.Services
}

//#
//# Specific methods
//#---------------------------------------------------------------------


//
//# AddNode: Add a new node into the cluster
func (c *Cluster) AddNode(n *ClusterNode) error {
  env.Output.WriteChDebug("(Cluster::AddNode) Add node '"+n.Name+"' to cluster")

  if c == nil {
      return errors.New("(Cluster::AddNode) Cluster not initialized")    
  }

  if c.Nodes == nil {
    env.Output.WriteChDebug("(Cluster::AddNode) Initializing cluster's Nodes")
    c.Nodes = make(map[string]*ClusterNode)
  }

  if _,exist := c.Nodes[n.Name]; exist {
    env.Output.WriteChWarn("(Cluster::AddNode) The node "+n.Name+" does already exist and will be overwritten.")
  }
  c.Nodes[n.Name] = n

  return nil
}

//#
//# Common methods
//#---------------------------------------------------------------------

//# String: convert a Cluster object to string
func (c *Cluster) String() string {
	return utils.ObjectToJsonString(c)
}

//#######################################################################################################
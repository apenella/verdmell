/*
Cluster Engine management

The package 'cluster' is used by verdmell to manage the cluster. 

=> Is known as a cluster a set of nodes.

-ClusterEngine
-Cluster
-ClusterNode
-ClusterService
-ClusterMessage
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
  Nodes map[string] *ClusterNode `json:"nodes"`
  //map to store the services that belong to the cluster
  Services map[string] *ClusterService `json:"services"` 
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
  return c.Nodes
}

//
//# SetServices: set attribute from Cluster
func (c *Cluster) SetServices(services map[string] *ClusterService) {
  env.Output.WriteChDebug("(Cluster::SetServices) Set Services' value")
  c.Services = services
}

//
//# GetServices: get attribute from Cluster
func (c *Cluster) GetServices() map[string] *ClusterService {
  env.Output.WriteChDebug("(Cluster::GetNodes) Get Services' value")
  return c.Services
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# GetNode: return a node from the cluster
func (c *Cluster) GetNode(name string) (error, *ClusterNode) {
  env.Output.WriteChDebug("(Cluster::GetNode) Retrieve node '"+name+"' from cluster")

  if node, exist := c.Nodes[name]; !exist {
    msg := "(Cluster::GetNode) Node '"+name+"' does not exit on the cluster"
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  } else {
    return nil, node
  }
}
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
    env.Output.WriteChWarn("(Cluster::AddNode) Node "+n.Name+" does already exist and will be overwritten.")
  }
  c.Nodes[n.Name] = n

  return nil
}
//
//# GetService: return a service from the cluster
func (c *Cluster) GetService(name string) (error, *ClusterService) {
  env.Output.WriteChDebug("(Cluster::GetService) Retrieve service '"+name+"' from cluster")

  if service, exist := c.Services[name]; !exist {
    msg := "(Cluster::GetService) Service '"+name+"' does not exit on the cluster"
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  } else {
    return nil, service
  }
}
//
//# AddNode: Add a new node into the cluster
func (c *Cluster) AddService(s *ClusterService) error {
  env.Output.WriteChDebug("(Cluster::AddService) Add service '"+s.Name+"' to cluster")

  if c == nil {
      return errors.New("(Cluster::AddService) Cluster not initialized")    
  }

  if c.Services == nil {
    env.Output.WriteChDebug("(Cluster::AddService) Initializing cluster's Nodes")
    c.Services = make(map[string]*ClusterService)
  }

  if _,exist := c.Services[s.Name]; exist {
    env.Output.WriteChWarn("(Cluster::AddService) Service "+s.Name+" does already exist and will be overwritten.")
  }
  c.Services[s.Name] = s

  return nil
}


//#
//# Common methods
//#---------------------------------------------------------------------

//# String: convert a Cluster object to string
func (c *Cluster) String() string {
  if err, str := utils.ObjectToJsonString(c); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//#######################################################################################################
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
  "strconv"
  "verdmell/service"
  "verdmell/utils"
)

//#
//#
//# ClusterService struct:
//# ClusterService is defined by a node name and its URL
type ClusterService struct{
	Name string `json:"name"`
  Status int `json:"status"`
	ServiceNodes map[string]*service.ServiceObject `json:"nodes"`
	CandidateForDelation bool `json:"candidatefordeletion"`
}
//
//# NewClusterNode: return a ClusterService instance to be run
func NewClusterService(name string) (error, *ClusterService) {
	env.Output.WriteChDebug("(ClusterService::NewClusterService)")
  srv := new(ClusterService)

  srv.SetName(name)
  srv.SetStatus(-1)
  srv.SetCandidateForDelation(false)

  srv.ServiceNodes = make( map[string] *service.ServiceObject )
	
  return nil, srv
}

//#
//# Getters and Setters
//#----------------------------------------------------------------------------------------

//
//# SetName: attribute from ClusterService
func (c *ClusterService) SetName(name string) {
  env.Output.WriteChDebug("(ClusterService::SetName) Set value '"+name+"'")
  c.Name = name
}
//
//# SetStatus: attribute from ClusterService
func (c *ClusterService) SetStatus(status int) {
  env.Output.WriteChDebug("(ClusterService::SetStatus) Set value")
  c.Status = status
}
//
//# SetServiceNodes: attribute from ClusterService
func (c *ClusterService) SetServiceNodes(s map[string] *service.ServiceObject) {
  env.Output.WriteChDebug("(ClusterService::SetServiceNodes) Set value")
  c.ServiceNodes = s
}
//
//# SetCandidateForDelation: attribute from ClusterNode
func (c *ClusterService) SetCandidateForDelation(d bool) {
  env.Output.WriteChDebug("(ClusterService::SetCandidateForDelation)")
  c.CandidateForDelation = d
}
//
//# GetName: attribute from ClusterService
func (c *ClusterService) GetName() string {
  env.Output.WriteChDebug("(ClusterService::GetName) Get value '"+c.Name+"'")
  return c.Name
}
//
//# GetStatus: attribute from ClusterService
func (c *ClusterService) GetStatus() int {
  env.Output.WriteChDebug("(ClusterService::SetStatus) Get value")
  return c.Status
}
//
//# GetServiceNodes: attribute from ClusterService
func (c *ClusterService) GetServiceNodes() map[string] *service.ServiceObject {
  env.Output.WriteChDebug("(ClusterService::GetServiceNodes) Get value")
  return c.ServiceNodes
}
//
//# SetCandidateForDelation: attribute from ClusterNode
func (c *ClusterService) GetCandidateForDelation() bool {
  env.Output.WriteChDebug("(ClusterService::GetCandidateForDelation)")
  return c.CandidateForDelation
}

//#
//# Specific methods
//#---------------------------------------------------------------------
//
//# AddService: method add a service for the current node
func (c *ClusterService) AddServiceNode(n string, s *service.ServiceObject) error {
  env.Output.WriteChDebug("(ClusterService::AddNode) Add node '"+n+"' to service '"+c.GetName()+"'["+strconv.Itoa(int(s.GetTimestamp()))+"]")
  c.ServiceNodes[n] = s
  c.CalculateStatusForClusterService()
	return nil
}
//
//# GetServiceNode: get service from specifc node
func (c *ClusterService) GetServiceNode(n string) (error, *service.ServiceObject) {
  env.Output.WriteChDebug("(ClusterService::GetServiceNode) Get service '"+c.GetName()+"' from node '"+n+"'")
  
  if srv, exist := c.ServiceNodes[n]; !exist {
    msg := "(ClusterService::GetServiceNode) Service '"+c.GetName()+"' not exit for node '"+n+"'"
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  } else {
    return nil, srv    
  }
}

//
//# CalculateStatusForClusterService: method sends a sample to the sample channel
func (c *ClusterService) CalculateStatusForClusterService() error {
  env.Output.WriteChDebug("(ClusterService::CalculateStatusForClusterService) Calculate ClusterService '"+c.GetName()+"' status")
  //Exit codes
  // OK: 0
  // WARN: 1
  // ERROR: 2
  // UNKNOWN: others (-1)
  //
  currentStatus := -1

  for _,service := range c.GetServiceNodes() {
    status := service.GetStatus()
    //exitStatus calculates the status
    if currentStatus < status {
      currentStatus = status
    }
  }
  env.Output.WriteChDebug("(ClusterService::CalculateStatusForClusterService) Status for '"+c.GetName()+"' has changed")

  c.SetStatus(currentStatus)

  return nil
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# string: convert a ClusterService object to string
func (c *ClusterService) String() string {
  if err, str := utils.ObjectToJsonString(c); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//#######################################################################################################
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
  "verdmell/service"
  "verdmell/utils"
)

//#
//#
//# ClusterNode struct:
//# ClusterNode is defined by a node name and its URL
type ClusterNode struct{
	Name string `json:"name"`
	URL string	`json:"URL"`
  Timestamp int64 `json:"timestamp"`
  Service *service.ServiceObject `json:"service"`
  NodeServices map[string] *service.ServiceObject `json:"services"`
  CandidateForDetelion bool `json:"candidatefordeletion"`
}
//
//# NewClusterNode: return a CheckEngine instance to be run
func NewClusterNode(name string, url string) (error, *ClusterNode) {
	env.Output.WriteChDebug("(ClusterNode::NewClusterNode)")
	node := new(ClusterNode)

	node.SetName(name)
	node.SetURL(url)
  node.SetTimestamp(0)
  node.SetCandidateForDelation(false)

  node.NodeServices = make( map[string]*service.ServiceObject )

	return nil, node
}

//#
//# Getters and Setters
//#----------------------------------------------------------------------------------------

//
//# SetName: attribute from ClusterNode
func (c *ClusterNode) SetName(name string) {
  env.Output.WriteChDebug("(ClusterNode::SetName) Set value '"+name+"'")
  c.Name = name
}
//
//# SetURL: attribute from ClusterNode
func (c *ClusterNode) SetURL(url string) {
  env.Output.WriteChDebug("(ClusterNode::SetURL) Set value '"+url+"'")
  c.URL = url
}
//
//# SetTimestamp: attribute from ClusterNode
func (c *ClusterNode) SetTimestamp(t int64) {
  env.Output.WriteChDebug("(ClusterNode::SetTimestamp)")
  c.Timestamp = t
}
//
//# SetService: attribute from ClusterNode
func (c *ClusterNode) SetService(s *service.ServiceObject) {
  env.Output.WriteChDebug("(ClusterNode::SetService)")
  c.Service = s
}
//
//# SetNodeServices: attribute from ClusterNode
func (c *ClusterNode) SetNodeServices(s map[string] *service.ServiceObject) {
  env.Output.WriteChDebug("(ClusterNode::SetNodeServices)")
  c.NodeServices = s
}
//
//# SetCandidateForDetelion: attribute from ClusterNode
func (c *ClusterNode) SetCandidateForDelation(d bool) {
  env.Output.WriteChDebug("(ClusterNode::SetCandidateForDelation)")
  c.CandidateForDetelion = d
}
//
//# GetName: attribute from ClusterNode
func (c *ClusterNode) GetName() string {
  env.Output.WriteChDebug("(ClusterNode::GetName) Get value '"+c.Name+"'")
  return c.Name
}
//
//# GetURL: attribute from ClusterNode
func (c *ClusterNode) GetURL() string {
  env.Output.WriteChDebug("(ClusterNode::GetName) Get value '"+c.URL+"'")
  return c.URL
}
//
//# GetTimestamp: attribute from ClusterNode
func (c *ClusterNode) GetTimestamp() int64 {
  env.Output.WriteChDebug("(ClusterNode::GetTimestamp)")
  return c.Timestamp
}
//
//# GetService: attribute from ClusterNode
func (c *ClusterNode) GetService() *service.ServiceObject {
  env.Output.WriteChDebug("(ClusterNode::GetService)")
  return c.Service
}
//
//# GetNodeServices: attribute from ClusterNode
func (c *ClusterNode) GetNodeServices() map[string] *service.ServiceObject {
  env.Output.WriteChDebug("(ClusterNode::GetNodeServices)")
  return c.NodeServices
}
//
//# SetCandidateForDetelion: attribute from ClusterNode
func (c *ClusterNode) GetCandidateForDelation() bool {
  env.Output.WriteChDebug("(ClusterNode::GetCandidateForDetelion)")
  return c.CandidateForDetelion
}
//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# GetNodeStatus: method sets the Status value for the ServiceObject
func (c *ClusterNode) GetNodeStatus() (error, int) {
  env.Output.WriteChDebug("(ClusterNode::GetNodeStatus)")
  if c.Service == nil {
    return errors.New("(ClusterNode::GetNodeStatus) No service for '"+c.GetName()+"' has been defined"),-1
  }

  return nil, c.Service.GetStatus()
}
//
//# HasService: method return if a service is defined on current node
func (c *ClusterNode) HasService(s string) (error, *service.ServiceObject) {
  env.Output.WriteChDebug("(ClusterNode::HasService) "+s)

  if srv, exist := c.NodeServices[s]; !exist {
    msg := "(ClusterNode::HasService) Service '"+s+"' does not exit on node "+c.GetName()
    env.Output.WriteChDebug(msg)
    return errors.New(msg), nil
  } else {
    return nil, srv    
  }

}
//
//# AddService: method add a service for the current node
func (c *ClusterNode) AddService(s *service.ServiceObject) error {
  env.Output.WriteChDebug("(ClusterNode::AddService) Add service '"+s.GetName()+"' on node '"+c.GetName()+"'")
  if err,_ := c.HasService(s.GetName()); err != nil {
    c.NodeServices[s.GetName()] = s
    env.Output.WriteChDebug("(ClusterNode::AddService) Service '"+s.GetName()+"' added on node '"+c.GetName()+"'")
  } else {
    return err
  }

  return nil
}
//
//# IncreaseTimestamp: method add a service for the current node
func (c *ClusterNode) IncreaseTimestamp() error {
  env.Output.WriteChDebug("(ClusterNode::IncreaseTimestamp) Node '"+c.GetName()+"'")
  c.Timestamp++
  return nil
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# string: convert a ClusterService object to string
func (c *ClusterNode) String() string {
  if err, str := utils.ObjectToJsonString(c); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//#######################################################################################################
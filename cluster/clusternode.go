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
//# ClusterNode struct:
//# ClusterNode is defined by a node name and its URL
type ClusterNode struct{
	Name string `json:"name"`
	URL string	`json:"URL"`
  Timestamp int64 `json:"timestamp"`
  Status int `json:"status"`
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
  node.SetStatus(-1)
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
//# SetStatus: attribute from ClusterNode
func (c *ClusterNode) SetStatus(status int) {
  env.Output.WriteChDebug("(ClusterNode::SetStatus) Set value")
  c.Status = status
}
//
//# SetTimestamp: attribute from ClusterNode
func (c *ClusterNode) SetTimestamp(t int64) {
  env.Output.WriteChDebug("(ClusterNode::SetTimestamp)")
  c.Timestamp = t
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
//# GetStatus: attribute from ClusterNode
func (c *ClusterNode) GetStatus() int {
  env.Output.WriteChDebug("(ClusterNode::GetStatus) Get value")
  return c.Status
}
//
//# GetTimestamp: attribute from ClusterNode
func (c *ClusterNode) GetTimestamp() int64 {
  env.Output.WriteChDebug("(ClusterNode::GetTimestamp)")
  return c.Timestamp
}
//
//# GetNodeServices: attribute from ClusterNode
func (c *ClusterNode) GetNodeServices() map[string] *service.ServiceObject {
  env.Output.WriteChDebug("(ClusterNode::GetNodeServices)")
  return c.NodeServices
}
//
//# SetCandidateForDetelion: attribute from ClusterNode
func (c *ClusterNode) GetCandidateForDeletion() bool {
  env.Output.WriteChDebug("(ClusterNode::GetCandidateForDetelion)")
  return c.CandidateForDetelion
}
//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# CopyClusterNode: method copies a cluster node
func (c *ClusterNode) CopyClusterNode() *ClusterNode {
  env.Output.WriteChDebug("(ClusterNode::CopyClusterNode)")
  if c == nil {
    return nil
  }
  
  node := new(ClusterNode)
  node.SetName(c.GetName())
  node.SetURL(c.GetURL())
  node.SetStatus(c.GetStatus())
  node.SetTimestamp(c.GetTimestamp())
  node.SetCandidateForDelation(c.GetCandidateForDeletion())
  node.SetNodeServices(c.GetNodeServices())

  return node
}
//
//# GetNodeStatus: method sets the Status value for the ServiceObject
// func (c *ClusterNode) GetNodeStatus() (error, int) {
//   env.Output.WriteChDebug("(ClusterNode::GetNodeStatus)")
//   var err error
//   var service *service.ServiceObject

//   if err, service = c.HasService(c.GetName()); err != nil {
//     return errors.New("(ClusterNode::GetNodeStatus) "+err.Error()),-1
//   }

//   return nil, service.GetStatus()
// }
//
//# HasService: method return if a service is defined on cluster node
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
//# AddService: method add a service to cluster node
func (c *ClusterNode) AddService(s *service.ServiceObject) error {
  env.Output.WriteChDebug("(ClusterNode::AddService) Add service '"+s.GetName()+"' on node '"+c.GetName()+"'["+strconv.Itoa(int(s.GetTimestamp()))+"]")
  if err,_ := c.HasService(s.GetName()); err != nil {
    c.NodeServices[s.GetName()] = s
    env.Output.WriteChDebug("(ClusterNode::AddService) Service '"+s.GetName()+"' added on node '"+c.GetName()+"'")
    if s.GetName() == c.GetName() {
      c.SetStatus(s.GetStatus())
      env.Output.WriteChDebug("(ClusterNode::AddService) Node '"+c.GetName()+"' has changed its status")
    }
  } else {
    return err
  }

  return nil
}
//
//# IncreaseTimestamp: method increade the timestamp to cluster node
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
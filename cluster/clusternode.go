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
)

//#
//#
//# ClusterNode struct:
//# ClusterNode is defined by a node name and its URL
type ClusterNode struct{
	Name string `json:"name"`
	URL string	`json:"URL"`
  Service *service.ServiceObject `json:"service"`
}
//
//# NewClusterNode: return a CheckEngine instance to be run
func NewClusterNode(name string, url string) (error, *ClusterNode) {
	env.Output.WriteChDebug("(ClusterNode::NewClusterNode)")
  	node := new(ClusterNode)

  	node.SetName(name)
  	node.SetURL(url)
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
//# SetService: attribute from ClusterNode
func (c *ClusterNode) SetService(s *service.ServiceObject) {
  env.Output.WriteChDebug("(ClusterNode::SetService)")
  c.Service = s
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
//# GetService: attribute from ClusterNode
func (c *ClusterNode) GetService() *service.ServiceObject {
  env.Output.WriteChDebug("(ClusterNode::GetName) Get value '"+c.URL+"'")
  return c.Service
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# GetNodeStatus: method sets the Status value for the ServiceObject
func (c *ClusterNode) GetNodeStatus() (error, int) {
  env.Output.WriteChDebug("(ClusterNode::GetNodeStatus)")
  if c.Service == nil {
    return errors.New("(ClusterNode::GetNodeStatus) The service for '"+c.GetName()+"' has not been defined"),-1
  }

  return nil, c.Service.GetStatus()
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# string: convert a ClusterEngine object to string
func (c *ClusterNode) String() string {
  return "'"+c.Name+"': '"+c.URL+"'"
}

//#######################################################################################################
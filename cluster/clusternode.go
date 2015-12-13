/*
Cluster Engine management

The package 'cluster' is used by verdmell to manage the cluster. 

=> Is known as a cluster a set of nodes.

-ClusterEngine
-Cluster
-ClusterNode
*/
package cluster

//#
//#
//# ClusterNode struct:
//# ClusterNode is defined by a node name and its URL
type ClusterNode struct{
	Name string `json:"name"`
	URL string	`json:"URL"` 
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
//# GetName: attribute from ClusterNode
func (c *ClusterNode) GetName() string{
  env.Output.WriteChDebug("(ClusterNode::GetName) Get value '"+c.Name+"'")
  return c.Name
}
//
//# GetURL: attribute from ClusterNode
func (c *ClusterNode) GetURL() string{
  env.Output.WriteChDebug("(ClusterNode::GetName) Get value '"+c.URL+"'")
  return c.URL
}

//#
//# Common methods
//#---------------------------------------------------------------------

//# string: convert a ClusterEngine object to string
func (c *ClusterNode) String() string {
  return "'"+c.Name+"': '"+c.URL+"'"
}

//#######################################################################################################
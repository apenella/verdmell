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
	"math"
	"strconv"
	"verdmell/utils"
)

//#
//#
//# Cluster struct:
//# Cluster is a set of nodes
type Cluster struct {
	// map to store the nodes that belong to the cluster
	Nodes map[string]*ClusterNode `json:"nodes"`
	// map to store the services that belong to the cluster
	Services map[string]*ClusterService `json:"services"`
	// map to control candidates for deletion
	candidatesForDeletion map[string]map[string]bool
}

//
//# NewCluster: return a CheckEngine instance to be run
func NewCluster() (error, *Cluster) {
	// func NewCluster(n map[string]*ClusterNode, s map[string]string) (error, *Cluster) {
	env.Output.WriteChDebug("(Cluster::NewCluster)")
	cluster := new(Cluster)
	cluster.candidatesForDeletion = make(map[string]map[string]bool)

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
func (c *Cluster) SetServices(services map[string]*ClusterService) {
	env.Output.WriteChDebug("(Cluster::SetServices) Set Services' value")
	c.Services = services
}

//
//# GetServices: get attribute from Cluster
func (c *Cluster) GetServices() map[string]*ClusterService {
	env.Output.WriteChDebug("(Cluster::GetNodes) Get Services' value")
	return c.Services
}

//
//# GetCandidatesForDeletion
func (c *Cluster) GetCandidatesForDeletion() map[string]map[string]bool {
	env.Output.WriteChDebug("(Cluster::GetCandidatesForDeletion) Get value")
	return c.candidatesForDeletion
}

//
//# SetCandidatesForDeletion
func (c *Cluster) SetCandidatesForDeletion(cs map[string]map[string]bool) {
	env.Output.WriteChDebug("(Cluster::SetCandidatesForDeletion) Set value")
	c.candidatesForDeletion = cs
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# GetNode: return a node from the cluster
func (c *Cluster) GetNode(name string) (error, *ClusterNode) {
	env.Output.WriteChDebug("(Cluster::GetNode) Retrieve node '" + name + "' from cluster")

	if node, exist := c.Nodes[name]; !exist {
		msg := "(Cluster::GetNode) Node '" + name + "' does not exit on the cluster"
		env.Output.WriteChDebug(msg)
		return errors.New(msg), nil
	} else {
		return nil, node
	}
}

//
//# AddNode: Add a new node into cluster
func (c *Cluster) AddNode(n *ClusterNode) error {
	env.Output.WriteChDebug("(Cluster::AddNode) Add node '" + n.Name + "' to cluster")
	// validate cluster
	if c == nil {
		return errors.New("(Cluster::AddNode) Cluster not initialized")
	}
	// validate if exist any node
	if c.Nodes == nil {
		env.Output.WriteChDebug("(Cluster::AddNode) Initializing cluster's Nodes")
		c.Nodes = make(map[string]*ClusterNode)
	}

	env.Output.WriteChDebug("(Cluster::AddNode) Node '" + n.GetName() + "' {status:" + strconv.Itoa(n.GetStatus()) + ", timestamp:" + strconv.Itoa(int(n.GetTimestamp())) + "}")
	c.Nodes[n.GetName()] = n

	return nil
}

//
//# DeleteNode: Delete node from cluster
func (c *Cluster) DeleteNode(node string) error {
	if c == nil {
		return errors.New("(Cluster::DeleteNode) Cluster not initialized")
	}
	delete(c.Nodes, node)
	return nil
}

//
//# GetService: return a service from the cluster
func (c *Cluster) GetService(name string) (error, *ClusterService) {
	env.Output.WriteChDebug("(Cluster::GetService) Retrieve service '" + name + "' from cluster")

	if service, exist := c.Services[name]; !exist {
		msg := "(Cluster::GetService) Service '" + name + "' does not exit on the cluster"
		env.Output.WriteChDebug(msg)
		return errors.New(msg), nil
	} else {
		return nil, service
	}
}

//
//# AddNode: Add a new node into the cluster
func (c *Cluster) AddService(s *ClusterService) error {
	env.Output.WriteChDebug("(Cluster::AddService) Add service '" + s.Name + "' to cluster")

	if c == nil {
		return errors.New("(Cluster::AddService) Cluster not initialized")
	}

	if c.Services == nil {
		env.Output.WriteChDebug("(Cluster::AddService) Initializing cluster's Nodes")
		c.Services = make(map[string]*ClusterService)
	}

	if _, exist := c.Services[s.Name]; exist {
		env.Output.WriteChWarn("(Cluster::AddService) Service " + s.Name + " does already exist and will be overwritten.")
	}
	c.Services[s.Name] = s

	return nil
}

//
//# DeleteService: Delete service from cluster
func (c *Cluster) DeleteService(service string) error {
	if c == nil {
		return errors.New("(Cluster::DeleteService) Cluster not initialized")
	}
	delete(c.Services, service)
	return nil
}

//
//# AddCandidatesForDeletion
func (c *Cluster) ConsensusForDeletion(candidate string, from string, reached bool) (error, bool) {
	env.Output.WriteChDebug("(Cluster::ConsensusForDeletion)")
	// consensus will determine how many nodes are required to determine if a node have to be deleted
	consensusNumber := int(math.Log2(float64(len(c.GetNodes()))))
	deletable := false

	// get information from candidate
	if exist, _ := c.candidatesForDeletion[candidate][from]; !exist {
		// if candidate is not reached then it should be subject for discussion
		if !reached {
			env.Output.WriteChDebug("(Cluster::ConsensusForDeletion) Node '" + candidate + "' is unreachable from '" + from + "'")
			if c.candidatesForDeletion[candidate] == nil {
				c.candidatesForDeletion[candidate] = make(map[string]bool)
			}
			c.candidatesForDeletion[candidate][from] = true
			// node is deletable when number of nodes which candidate is unreachable is over consensus
			if len(c.candidatesForDeletion[candidate]) > consensusNumber {
				env.Output.WriteChDebug("(Cluster::ConsensusForDeletion) Node '" + candidate + "' will be deleted'")
				deletable = true
			}
		}
	} else {
		// candidate already exist
		if reached {
			env.Output.WriteChDebug("(Cluster::ConsensusForDeletion) Node '" + candidate + "' is already reachable from '" + from + "'")
			// should be delete node as an unreached from
			delete(c.candidatesForDeletion[candidate], from)
		}
	}
	return nil, deletable
}

//
// String method convert a Cluster object to string
func (c *Cluster) String() string {
	var str string
	var err error

	str, err = utils.ObjectToJSONString(c)
	if err != nil {
		return err.Error()
	}

	return str
}

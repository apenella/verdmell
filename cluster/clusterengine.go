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
	"verdmell/environment"
  "verdmell/service"
  "verdmell/utils"
)

//
var env *environment.Environment

//#
//#
//# ClusterEngine struct:
//# ClusterEngine defines a map to store the maps
type ClusterEngine struct{
	//Ui ui.UI
  Cluster *Cluster `json:"cluster"`
  inputChannel chan interface{} `json:"-"`
  outputChannels map[chan []byte] bool `json:"-"`
}
//
//# NewClusterEngine: return a CheckEngine instance to be run
func NewClusterEngine(e *environment.Environment) (error, *ClusterEngine){
  e.Output.WriteChDebug("(ClusterEngine::NewClusterEngine)")
  
  cluster := new(ClusterEngine)

  // Get the environment attributes
  env = e

  // Set Cluster to engine, where will be defined nodes and services
  if err, c := NewCluster(); err != nil {
    return err, nil
  } else {
    cluster.SetCluster(c)
  }

  // Check is user has defined any node
  if env.Setup.Cluster != nil {
    env.Output.WriteChDebug("(ClusterEngine::NewClusterEngine) There are some nodes defined by user.")

    // each node defined by user is going to be add on cluster 
    for name, url := range env.Setup.Cluster {
      if err ,node := NewClusterNode(name,url); err != nil {
        return err, nil
      } else {
        if err = cluster.AddNode(node); err != nil {
          return err, nil
        }
        env.Output.WriteChDebug("(ClusterEngine::NewClusterEngine) The node '"+name+"' has been added into the cluster")
        env.Output.WriteChDebug("(ClusterEngine::NewClusterEngine) The cluster nodes '"+cluster.Cluster.String()+"'")
      }
    } 
  }

  // Initialize the outputSampleChannels
  cluster.outputChannels = make(map[chan []byte] bool)

  // Starting the method to receive data to handle
  if err := cluster.StartReceiver(); err != nil {
    return errors.New("(ClusterEngine::NewClusterEngine) The service receiver for could not be started"), nil  
  }

  // the cluster engine is set to environment
  env.SetClusterEngine(cluster)
  env.Output.WriteChDebug("(ClusterEngine::NewClusterEngine) I'm your new Cluster Engine instance.")

  return nil, cluster
}

//
//# SetCluster: attribute from Cluster
func (c *ClusterEngine) SetCluster(cluster *Cluster) {
  env.Output.WriteChDebug("(ClusterEngine::SetCluster) Set value '"+cluster.String()+"'")
  c.Cluster = cluster
}
//
//# SetInputChannel: attribute from Cluster
func (c *ClusterEngine) SetInputChannel(i chan interface{}) {
  env.Output.WriteChDebug("(ClusterEngine::SetInputChannel)")
  c.inputChannel = i
}
//
//# SetOutputChannel: attribute from Cluster
func (c *ClusterEngine) SetOutputChannels(o map[chan []byte] bool) {
  env.Output.WriteChDebug("(ClusterEngine::SetOutputChannels)")
  c.outputChannels = o
}
//
//# GetCluster: attribute from ClusterNode
func (c *ClusterEngine) GetCluster() *Cluster {
  env.Output.WriteChDebug("(ClusterEngine::GetCluster) Get value '"+c.Cluster.String()+"'")
  return c.Cluster
}
//
//# GetInputChannel: attribute from Cluster
func (c *ClusterEngine) GetInputChannel() chan interface{} {
  env.Output.WriteChDebug("(ClusterEngine::GetInputChannel)")
  return c.inputChannel
}
//
//# SetOutputChannel: attribute from Cluster
func (c *ClusterEngine) GetOutputChannels() map[chan []byte] bool {
  env.Output.WriteChDebug("(ClusterEngine::GetOutputChannels)")
  return c.outputChannels
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# SayHi: do nothing
func (c *ClusterEngine) SayHi() {
  env.Output.WriteChInfo("(ClusterEngine::SayHi) Hi! I'm your new Cluster Engine instance.")
}
//
//# AddNode: Add a new node into the cluster
func (c *ClusterEngine) AddNode(n *ClusterNode) error {
  env.Output.WriteChDebug("(ClusterEngine::AddNode) Add node '"+n.Name+"' to cluster")
  
  if c.Cluster == nil {
    env.Output.WriteChDebug("(ClusterEngine::AddNode) Initialize a Cluster instance")    
    _,c.Cluster = NewCluster()
  }

  return c.Cluster.AddNode(n)
}
//
//# StartServiceEngine: method prepares the system to wait sample and calculate the results for services
func (c *ClusterEngine) StartReceiver() error {
  c.inputChannel = make(chan interface{})

  env.Output.WriteChDebug("(ClusterEngine::StartReceiver) Starting services receiver")
  go func() {
    defer close (c.inputChannel)
    for{
      select{

      case unknownObject := <-c.inputChannel:
        env.Output.WriteChDebug("(ClusterEngine::StartReceiver) Received")
        switch object := unknownObject.(type){
        // a ServiceObject is received from the own server
        case *service.ServiceObject:
          c.handleServiceObject(object)
        // a []byte should be received from other cluster node
        case []byte:
          c.handleClusterMessage(object)
        }
      }
    }
  }()
  return nil
}
//
//# handleServiceObject: handles the incoming serviceObject messages, that messages are sent by own node
func (c *ClusterEngine) handleServiceObject(s *service.ServiceObject) {
  cluster := c.GetCluster()
  timestamp := s.GetTimestamp()

  env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) ServiceObject received. Service '"+s.GetName()+"' with timestamp:"+strconv.Itoa(int(timestamp)))

  // The current service is the node service
  if env.Setup.Hostname == s.GetName() {
    if err, node := cluster.GetNode(s.GetName()); err == nil {
      env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Current node's status received ")

      if node.GetService() == nil || timestamp >= node.GetService().GetTimestamp() {
        // Set the node status
        node.SetService(s)

        // Prepare the message to be sent to the cluster
        if err, message := NewClusterMessage(s.GetName(),timestamp,s); err == nil {
          env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) "+message.String())
          if err, messageBytes := utils.InterfaceToBytes(message); err == nil {
            c.sendData(messageBytes)
          } else {
            env.Output.WriteChError("(ClusterEngine::handleServiceObject) Error on sendData ("+err.Error()+")")
          }
        }else{  
          env.Output.WriteChError("(ClusterEngine::handleServiceObject) Error on NewClusterMessage ("+err.Error()+")")
        }
      }        
    }
  } else {
    env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Service '"+s.GetName()+"' status received")
  }
}
//
//# handleServiceObject: handles the incoming serviceObject messages
func (c *ClusterEngine) handleClusterMessage(data []byte) {
  env.Output.WriteChDebug("(ClusterEngine::handleClusterMessage) []byte received")
  //TODO
}
//
//# sendServicesStatus: method that send services to other engines
func (c *ClusterEngine) sendData(data []byte) error {
  env.Output.WriteChDebug("(ClusterEngine::sendData)")
  for c,_ := range c.GetOutputChannels(){
     env.Output.WriteChDebug("(ClusterEngine::sendData) Writing data on channel")
      c <- data
  }

  return nil
}
//
//# SendSample: method prepares the system to wait sample and calculate the results for services
func (c *ClusterEngine) SendService(service *service.ServiceObject) {
  env.Output.WriteChDebug("(ClusterEngine::SendService) Send service status for '"+service.GetName()+"'")
  c.inputChannel <- service
}
//
//# AddOutputSampleChan: Add a new channel to write service status
func (c *ClusterEngine) AddOutputChannel(o chan []byte) error {
  env.Output.WriteChDebug("(ClusterEngine::AddOutputChannel)")

  channels := c.GetOutputChannels()
  if _, exist := channels[o]; !exist {
    channels[o] = true
  } else {
    return errors.New("(ClusterEngine::AddOutputChannel) You are trying to add an existing channel")
  }

  return nil
}

//
//# GetCluster: get whole information from cluster
func (c *ClusterEngine) GetClusterInfo() (error,[]byte) {
  env.Output.WriteChDebug("(ClusterEngine::GetClusterData)")
  var cluster *Cluster

  if cluster = c.GetCluster(); cluster == nil {
    return errors.New("(ClusterEngine::GetClusterData) The cluster object has not been initialized"),nil
  }

  return nil,utils.ObjectToJsonByte(c.GetCluster())
}
//
//# GetClusterNodes: get all nodes from cluster
func (c *ClusterEngine) GetClusterNodes() (error,[]byte) {
  env.Output.WriteChDebug("(ClusterEngine::GetClusterNodes)")
  var cluster *Cluster

  if cluster = c.GetCluster(); cluster == nil {
    return errors.New("(ClusterEngine::GetClusterNodes) The cluster object has not been initialized"),nil
  }

  return nil,utils.ObjectToJsonByte(cluster.GetNodes())
}

//
// Gossip


//#
//# Common methods
//#---------------------------------------------------------------------

//# String: convert a ClusterEngine object to string
func (c *ClusterEngine) String() string {
  if err, str := utils.ObjectToJsonString(c.GetCluster()); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//#######################################################################################################
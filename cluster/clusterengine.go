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
  "bytes"
  "errors"
//  "strconv"
  "net/http"
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
  inputGossipChannel chan []byte `json:"-"`
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
        env.Output.WriteChDebug("(ClusterEngine::NewClusterEngine) Node '"+name+"' has been added into the cluster")
      }
    }
    env.Output.WriteChDebug("(ClusterEngine::NewClusterEngine) Cluster nodes '"+cluster.Cluster.String()+"'") 
  }

  // Initialize the outputSampleChannels
  cluster.outputChannels = make(map[chan []byte] bool)

  // Starting the method to receive data to handle
  if err := cluster.StartReceiver(); err != nil {
    return errors.New("(ClusterEngine::NewClusterEngine) Service receiver for could not be started"), nil  
  }

  // Initialize Gossip 
  cluster.inputGossipChannel = make(chan []byte)
  cluster.AddOutputChannel(cluster.inputGossipChannel)
  cluster.StartClusterGossip()

  // the cluster engine is set to environment
  env.SetClusterEngine(cluster)
  env.Output.WriteChInfo("(ClusterEngine::NewClusterEngine) I'm your new Cluster Engine instance.")

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
  env.Output.WriteChDebug("(ClusterEngine::GetCluster) Get value")
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
  env.Output.WriteChInfo("(ClusterEngine::SayHi) Hi! I'm your Cluster Engine instance.")
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
//# AddNode: Add a new node into the cluster
func (c *ClusterEngine) AddService(s *ClusterService) error {
  env.Output.WriteChDebug("(ClusterEngine::AddService) Add service '"+s.Name+"' to cluster")
  
  if c.Cluster == nil {
    env.Output.WriteChDebug("(ClusterEngine::AddService) Initialize a Cluster instance")    
    _,c.Cluster = NewCluster()
  }

  return c.Cluster.AddService(s)
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
          // call handleServiceObject for current node
          if err := c.handleServiceObject(object); err != nil {
            env.Output.WriteChDebug("(ClusterEngine::StartReceiver) "+err.Error())
          }
        // a []byte should be received from other cluster node
        case []byte:
          if err := c.handleClusterMessage(object); err != nil {
            env.Output.WriteChDebug("(ClusterEngine::StartReceiver) "+err.Error()) 
          }
        }
      }
    }
  }()
  return nil
}
//
//# handleServiceObject: handles the incoming serviceObject messages, that messages are sent by own node
func (c *ClusterEngine) handleServiceObject(s *service.ServiceObject) error {
  env.Output.WriteChDebug("(ClusterEngine::handleServiceObject)")
  var clusternode *ClusterNode
  var clusterservice *ClusterService
  nodename := env.Setup.Hostname

  if err, cluster := NewCluster(); err != nil {
    env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) "+err.Error())
  } else {
    if err, clusternode = c.Cluster.GetNode(nodename); err != nil {
      msg := "(ClusterEngine::handleServiceObject) "+err.Error()
      env.Output.WriteChDebug(msg)
      return errors.New(msg)
    }
    clusternode.AddService(s)
    cluster.AddNode(clusternode)


    if err, clusterservice = NewClusterService(s.GetName()); err != nil {
      msg := "(ClusterEngine::handleServiceObject) "+err.Error()
      env.Output.WriteChDebug(msg)
      return errors.New(msg)
    }
    clusterservice.AddServiceNode(nodename,s)
    cluster.AddService(clusterservice)

    env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Cluster prepared to be updated with service '"+s.GetName()+"'")

    if err := c.updateCluster(cluster); err != nil {
      env.Output.WriteChError("(ClusterEngine::handleServiceObject) "+err.Error())
    }
  }
  return nil
}

// func (c *ClusterEngine) handleServiceObject(s *service.ServiceObject) error {
//   env.Output.WriteChDebug("(ClusterEngine::handleServiceObject)")

//   var err error
//   var clusternode *ClusterNode
//   var clusterservice *ClusterService

//   cluster := c.GetCluster()
//   timestamp := s.GetTimestamp()
//   numberOfChecks := len(s.GetChecks())

//   //
//   // update cuurent clusternode
//   // current clusternode has been added on clusterengine creation
//   if err, clusternode = cluster.GetNode(env.Setup.Hostname); err != nil {
//     msg := "(ClusterEngine::handleServiceObject) "+err.Error()
//     env.Output.WriteChDebug(msg)
//     return errors.New(msg)
//   }

//   if err, clusterservice = cluster.GetService(s.GetName()); err != nil {
//     env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Create new service '"+s.GetName()+"' on cluster node '"+clusternode.GetName()+"'")
//     if err, clusterservice = NewClusterService(s.GetName()); err != nil {
//       msg := "(ClusterEngine::handleServiceObject) "+err.Error()
//       env.Output.WriteChDebug(msg)
//       return errors.New(msg)
//     }
//     c.AddService(clusterservice)
//   }

//   env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) ServiceObject received[node:"+clusternode.GetName()+", service:"+s.GetName()+", service timestamp:"+strconv.Itoa(int(timestamp))+", numchecks:"+strconv.Itoa(numberOfChecks)+"]")
//   //
//   // if clusternode has s service update its status
//   if err, service := clusternode.HasService(s.GetName()); err != nil {
//     env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Add new service '"+s.GetName()+"' on cluster node '"+clusternode.GetName()+"'")
//     c.updateCluster(clusternode,clusterservice,env.Setup.Hostname,s)
//   } else {
//     // clusternode doesn't already have s service
//     if timestamp >= service.GetTimestamp() ||  int(timestamp) <= numberOfChecks { 
//       //add service to node
//       env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Update service '"+service.GetName()+"' on cluster node '"+clusternode.GetName()+"'")
//       c.updateCluster(clusternode,clusterservice,env.Setup.Hostname,s)
//     }  
//   }
//   return nil
// }
//
//# handleServiceObject: handles the incoming serviceObject messages
func (c *ClusterEngine) handleClusterMessage(data []byte) error {
  env.Output.WriteChDebug("(ClusterEngine::handleClusterMessage)")
   
  //
  // Decoding received []byte to a cluster message
  if err, message := DecodeClusterMessage(data); err != nil {
    // When the data could not be decoded an error is thrown
    return errors.New("(ClusterMessage::handleClusterMessage) "+err.Error())
  } else {
    env.Output.WriteChDebug("(ClusterEngine::handleClusterMessage) Message from '"+message.GetFrom()+"'")
    // Decoding data from message
    if err, messageData := DecodeData(message.GetData()); err != nil {
      // When data coul not be decoded an error is thrown
      msg := "(ClusterEngine::handleClusterMessage) "+err.Error()
      env.Output.WriteChError(msg)
      return errors.New(msg)
    } else {
      cluster := c.GetCluster()
      //
      // Read data from the message: 
      // get the reviced message timestamp
      messageTimestamp := message.GetTimestamp()
      // get sender clusternode's stored data
      if err, clusternode := cluster.GetNode(message.GetFrom()); err != nil {
        env.Output.WriteChDebug("(ClusterEngine::handleClusterMessage) Received a message from new node '"+message.GetFrom()+"', and will be added on cluster")
        //
        // update information from nodes
        if err := c.updateCluster(messageData); err != nil {
          env.Output.WriteChError("(ClusterEngine::handleClusterMessage) "+err.Error())
        }
        //
        //update services
      } else {
        env.Output.WriteChDebug("(ClusterEngine::handleClusterMessage) Node '"+message.GetFrom()+"' has data on cluster")
        // get received node from stored timestamp
        nodeTimestamp := clusternode.GetTimestamp()
        //
        // if message timestamp is greater than nodeTimestamp means that the information from this node is newer
        if messageTimestamp > nodeTimestamp {
          // update information from nodes
          if err := c.updateCluster(messageData); err != nil {
            env.Output.WriteChError("(ClusterEngine::handleClusterMessage) "+err.Error())
          }
        }
      }
    }
  }
  return nil
}
//
//# updateCluster: method to update stored nodes information
func (c *ClusterEngine) updateCluster(clusterFrom *Cluster) error {
  env.Output.WriteChDebug("(ClusterEngine::updateCluster)")
  modified := false
  cluster := c.GetCluster()


  for _,nodeFrom := range clusterFrom.GetNodes(){
    // get stored node's data from cluster
    env.Output.WriteChDebug("(ClusterEngine::updateCluster) Received data from node '"+nodeFrom.GetName()+"' will be analized")
    if err, nodeStored := cluster.GetNode(nodeFrom.GetName()); err != nil {
      env.Output.WriteChWarn("(ClusterEngine::updateCluster) Node '"+nodeFrom.GetName()+"' will be added as controled node")
      // add the node to cluster
      cluster.AddNode(nodeFrom)
      //
      // for each service defined on node, check it's status
      for _,serviceFrom := range nodeFrom.GetNodeServices(){
        // get cluster service stored for serviceFrom
        if err, clusterService := cluster.GetService(serviceFrom.GetName()); err != nil {
          env.Output.WriteChWarn("(ClusterEngine::updateCluster) Service '"+serviceFrom.GetName()+"' will be added as controled service")
          // Create cluster service
          if err, clusterservice := NewClusterService(serviceFrom.GetName()); err != nil {
            env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
          } else {
            c.AddService(clusterservice)
            clusterservice.AddServiceNode(nodeFrom.GetName(),serviceFrom)
          }
        } else {
          // get stored data for the cluster service
          if err, serviceStored := clusterService.GetServiceNode(nodeFrom.GetName()); err != nil{
            env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
          } else {
            // update service whem serviceFrom timestamp is greater than the stored one
            if serviceFrom.GetTimestamp() > serviceStored.GetTimestamp() {
              clusterService.AddServiceNode(nodeFrom.GetName(),serviceFrom)
            }
          }
        }
      }
      modified = true
    } else {
      // update stored node when received timestamp from node is greater
      env.Output.WriteChDebug("(ClusterEngine::updateCluster) Node '"+nodeFrom.GetName()+"' already exist on cluster")
      if err, nodeFromService := nodeFrom.HasService(nodeFrom.GetName()); err != nil {
        env.Output.WriteChDebug("(ClusterEngine::updateCluster) Could not load service from '"+nodeFrom.GetName()+"'."+err.Error())
      } else {
        // get number of checks from the service. Is possible to detect a restart when number of checks is lower than timestamp
        numberOfChecks := len(nodeFromService.GetChecks())
        if ( nodeFrom.GetTimestamp() > nodeStored.GetTimestamp()) || int(nodeFrom.GetTimestamp()) <= numberOfChecks {
          env.Output.WriteChDebug("(ClusterEngine::updateCluster) Node '"+nodeFrom.GetName()+"' will be updated")
          cluster.AddNode(nodeFrom)
          //
          // for each service defined on node, check it's status
          for _,serviceFrom := range nodeFrom.GetNodeServices(){
            // get cluster service stored for serviceFrom
            if err, clusterService := cluster.GetService(serviceFrom.GetName()); err != nil {
              env.Output.WriteChWarn("(ClusterEngine::updateCluster) Service '"+serviceFrom.GetName()+"' will be added as controled service")
              // Create cluster service
              if err, clusterservice := NewClusterService(serviceFrom.GetName()); err != nil {
                env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
              } else {
                c.AddService(clusterservice)
                clusterservice.AddServiceNode(nodeFrom.GetName(),serviceFrom)
              }
            } else {
              // get stored data for the cluster service
              if err, serviceStored := clusterService.GetServiceNode(nodeFrom.GetName()); err != nil{
                env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
              } else {
                // update service whem serviceFrom timestamp is greater than the stored one
                if serviceFrom.GetTimestamp() > serviceStored.GetTimestamp() {
                  clusterService.AddServiceNode(nodeFrom.GetName(),serviceFrom)
                }            
              }
            }
          }
          modified = true
        }
      }
    }
  }

  // On changes, increase timestamp and deploy clusterMessage to cluster
  if modified {
    if err, clusternode := cluster.GetNode(env.Setup.Hostname); err != nil {
      env.Output.WriteChError("(ClusterEngine::updateCluster) Cluster status has changed but ClusterMessage could not be deploy due an error: "+err.Error())
    } else {
      clusternode.IncreaseTimestamp()

     // Prepare the message to be sent to the cluster
      if err, message := NewClusterMessage(env.Setup.Hostname, clusternode.GetTimestamp(),cluster); err == nil {
        env.Output.WriteChDebug("(ClusterEngine::updateCluster) New message ready to be sent")
        // transform data as []byte before send it
        if err, messageBytes := utils.InterfaceToBytes(message); err == nil {
          // send data to the outputChannels
          c.SendData(messageBytes)
        } else {
          env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
        }    
      }
    }
  }
  return nil
}
//updateservice status

//
//# updateCluster: method set new values to node and service
// func (c *ClusterEngine) updateCluster(clusternode *ClusterNode, clusterservice *ClusterService, n string, s *service.ServiceObject) error {
//   env.Output.WriteChDebug("(ClusterEngine::updateCluster) Update node '"+n+"' and service '"+s.GetName()+"'")
//   cluster := c.GetCluster()

//   clusternode.AddService(s)
//   clusterservice.AddServiceNode(n,s)
//   // any change changes clusternode timestamp
//   clusternode.IncreaseTimestamp()

//   // Prepare the message to be sent to the cluster
//   if err, message := NewClusterMessage(n,clusternode.GetTimestamp(),cluster); err == nil {
//     env.Output.WriteChDebug("(ClusterEngine::updateCluster) New message ready to be sent")
//     // transform data as []byte before send it
//     if err, messageBytes := utils.InterfaceToBytes(message); err == nil {
//       // send data to the outputChannels
//       c.SendData(messageBytes)
//     } else {
//       return errors.New("(ClusterEngine::updateCluster) "+err.Error())
//     }    
//   }
//   return nil
// }
//
//# SendData: method that send services to other engines
func (c *ClusterEngine) SendData(data []byte) error {
  env.Output.WriteChDebug("(ClusterEngine::SendData)")
  for c,_ := range c.GetOutputChannels(){
    env.Output.WriteChDebug("(ClusterEngine::SendData) Writing data on channel")
    c <- data
  }
  return nil
}
//
//# ReceiveData: method to receive data from the outdoor
func (c *ClusterEngine) ReceiveData(data interface{} ) {
  env.Output.WriteChDebug("(ClusterEngine::ReceiveData) Data received to be queued")
  c.inputChannel <- data
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
    return errors.New("(ClusterEngine::GetClusterData) Cluster object has not been initialized"),nil
  }

  return nil,utils.ObjectToJsonByte(c.GetCluster())
}
//
//# GetClusterNodes: get all nodes from cluster
func (c *ClusterEngine) GetClusterNodes() (error,[]byte) {
  env.Output.WriteChDebug("(ClusterEngine::GetClusterNodes)")
  var cluster *Cluster

  if cluster = c.GetCluster(); cluster == nil {
    return errors.New("(ClusterEngine::GetClusterNodes) Cluster object has not been initialized"),nil
  }

  return nil,utils.ObjectToJsonByte(cluster.GetNodes())
}
//
//# GetClusterServices: get all nodes from cluster
func (c *ClusterEngine) GetClusterServices() (error,[]byte) {
  env.Output.WriteChDebug("(ClusterEngine::GetClusterServices)")
  var cluster *Cluster

  if cluster = c.GetCluster(); cluster == nil {
    return errors.New("(ClusterEngine::GetClusterServices) Cluster object has not been initialized"),nil
  }

  return nil,utils.ObjectToJsonByte(cluster.GetServices())
}

//
// Gossip
//-----------------------------------------------------------------------

//# SelectNodesToGossip
func (c *ClusterEngine) SelectNodesToGossip() (error, map[string]*ClusterNode) {
  cluster := c.GetCluster()
  return nil,cluster.GetNodes()
}

//
//# StartClusterGossip
func (c *ClusterEngine) StartClusterGossip() error {
  env.Output.WriteChDebug("(ClusterEngine::StartClusterGossip) Starting cluster gossip")
  
  
  go func() {
    defer close (c.inputGossipChannel)
    for{
      select{
      case message := <-c.inputGossipChannel:
        env.Output.WriteChDebug("(ClusterEngine::StartClusterGossip) New message to gossip")
        if err, nodes := c.SelectNodesToGossip(); err == nil {
          for _, clusternode := range nodes {
            go func() {
              env.Output.WriteChDebug("(ClusterEngine::StartClusterGossip) Send message to node '"+clusternode.GetName()+"'")
              if err := c.SendGossipMessage(clusternode.GetURL(),message); err != nil {
                env.Output.WriteChDebug("(ClusterEngine::StartClusterGossip) "+err.Error())
              }
            }()
          }
        }
      }
    }
  }()
  return nil
}
//
//# SendGossipMessage
func (c *ClusterEngine) SendGossipMessage(url string, message []byte) error {
  
  uri := "/api/cluster"
  url = url + uri
  env.Output.WriteChDebug("(ClusterEngine::SendGossipMessage) Message to '"+url+"'")

  req, err := http.NewRequest("PUT", url, bytes.NewBuffer(message))
  req.Header.Set("X-Verdmell-From", env.Setup.Hostname)
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)


  if err != nil {
    msg := "(ClusterEngine::SendGossipMessage) "+err.Error()
    env.Output.WriteChError(msg)
    return errors.New(msg)
  }
  defer resp.Body.Close()

  return nil
}

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
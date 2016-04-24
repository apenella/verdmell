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
  "strconv"
  "strings"
  "net/http"
  "math"
  "math/rand"
  "time"
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
  outputChannels map[chan []byte] string `json:"-"`
  inputSyncChannel chan []byte `json:"-"`
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
  cluster.outputChannels = make(map[chan []byte] string)

  // Starting the method to receive data to handle
  if err := cluster.StartReceiver(); err != nil {
    return errors.New("(ClusterEngine::NewClusterEngine) Service receiver for could not be started"), nil  
  }

  // Initialize Sync 
  cluster.inputSyncChannel = make(chan []byte)
  cluster.AddOutputChannel(cluster.inputSyncChannel,"inputSyncChannel")
  cluster.StartClusterSync()

  // the cluster engine is set to environment
  env.SetClusterEngine(cluster)
  env.Output.WriteChInfo("(ClusterEngine::NewClusterEngine) I'm your new Cluster Engine instance.")

  return nil, cluster
}

//#
//# Getters/Setters methods for Checks object
//#---------------------------------------------------------------------

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
func (c *ClusterEngine) SetOutputChannels(o map[chan []byte] string) {
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
func (c *ClusterEngine) GetOutputChannels() map[chan []byte] string {
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
//# AddNode: Add a new node into cluster
func (c *ClusterEngine) AddNode(n *ClusterNode) error {
  env.Output.WriteChDebug("(ClusterEngine::AddNode) Add node '"+n.Name+"' to cluster")
  
  if c.Cluster == nil {
    env.Output.WriteChDebug("(ClusterEngine::AddNode) Initialize a Cluster instance")    
    _,c.Cluster = NewCluster()
  }

  return c.Cluster.AddNode(n)
}
//
//# DeleteNode: delete node from cluster
func (c *ClusterEngine) DeleteNode(nodename string) error {
  env.Output.WriteChDebug("(ClusterEngine::AddNode) DeleteNode node '"+nodename+"' from cluster")
  var node *ClusterNode
  var err error

  if c.Cluster == nil {
    return errors.New("(ClusterEngine::DeleteNode) Cluster has not been initialized")
  }
  if err, node = c.Cluster.GetNode(nodename); err != nil {
    return errors.New("(ClusterEngine::DeleteNode) "+err.Error())
  }
  // delete node from cluster
  c.Cluster.DeleteNode(nodename)
  // deassign node from services
  for servicename,_ := range node.GetNodeServices() {
    _, service := c.Cluster.GetService(servicename)
    service.DeleteServiceNode(nodename)
    // delete service from cluster when service has no node 
    if service.CountServiceNodes() < 1 {
      c.Cluster.DeleteService(servicename)
    }
  }
  return nil
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
        switch object := unknownObject.(type){
        // a ServiceObject is received from the own server
        case *service.ServiceObject:
          // call handleServiceObject for current node
          env.Output.WriteChDebug("(ClusterEngine::StartReceiver) Received serviceObject {service:'"+object.GetName()+", service_status:"+strconv.Itoa(object.GetStatus())+", service_timestamp:"+strconv.Itoa(int(object.GetTimestamp()))+"}")
          if err := c.handleServiceObject(object); err != nil {
            env.Output.WriteChError("(ClusterEngine::StartReceiver) "+err.Error())
          }
        // a []byte should be received from other cluster node
        case []byte:
          env.Output.WriteChDebug("(ClusterEngine::StartReceiver) Received []byte")
          if err := c.handleClusterMessage(object); err != nil {
            env.Output.WriteChError("(ClusterEngine::StartReceiver) "+err.Error()) 
          }
        }
      }
    }
  }()
  return nil
}
//
//# SendData: method that send services to other engines
func (c *ClusterEngine) SendData(data []byte) error {
  env.Output.WriteChDebug("(ClusterEngine::SendData)")
  for c, desc := range c.GetOutputChannels(){
    env.Output.WriteChDebug("(ClusterEngine::SendData) Writing data on channel '"+desc+"'")
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
func (c *ClusterEngine) AddOutputChannel(o chan []byte, desc string) error {
  env.Output.WriteChDebug("(ClusterEngine::AddOutputChannel) ")

  if o == nil {
    msg := "(ClusterEngine::AddOutputChannel) Null outputChannel"
    env.Output.WriteChError(msg)
    return errors.New(msg)
  }

  channels := c.GetOutputChannels()
  if _, exist := channels[o]; !exist {
    env.Output.WriteChDebug("(ClusterEngine::AddOutputChannel) New outputChannel registered")    
    channels[o] = desc
  } else {
    return errors.New("(ClusterEngine::AddOutputChannel) You are trying to add an existing channel")
  }
  return nil
}
//
//# handleServiceObject: handles the incoming serviceObject messages, that messages are sent by own node
func (c *ClusterEngine) handleServiceObject(s *service.ServiceObject) error {
  env.Output.WriteChDebug("(ClusterEngine::handleServiceObject)")
  var clusternode *ClusterNode
  var newclusternode *ClusterNode
  var clusterservice *ClusterService
  modify := false
  nodename := env.Whoami()

  if err, cluster := NewCluster(); err != nil {
    env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) "+err.Error())
  } else {
    if err, clusternode = c.Cluster.GetNode(nodename); err != nil {
      msg := "(ClusterEngine::handleServiceObject) "+err.Error()
      env.Output.WriteChDebug(msg)
      return errors.New(msg)
    }
    // to avoid traffic is check timestamp
    if err, nodeFromService := clusternode.HasService(s.GetName()); err != nil {
      env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Service will be assigned to node {node:'"+nodename+"', service:'"+s.GetName()+"', service_status:"+strconv.Itoa(int(s.GetStatus())) +", service_timestamp:"+strconv.Itoa(int(s.GetTimestamp()))+"}")
      modify = true
    } else {
      env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) '"+nodeFromService.GetName()+"' serviceStored {service_status:"+strconv.Itoa(nodeFromService.GetStatus())+", service_timestamp:"+strconv.Itoa(int(nodeFromService.GetTimestamp()))+"}")
      env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) '"+s.GetName()+"' serviceFrom {service_status:"+strconv.Itoa(s.GetStatus())+", service_timestamp:"+strconv.Itoa(int(s.GetTimestamp()))+"}")
      if s.GetTimestamp() > nodeFromService.GetTimestamp() {
        env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Service will be modified on node {node:'"+nodename+"', service:'"+s.GetName()+"', service_status:"+strconv.Itoa(int(s.GetStatus())) +", service_timestamp:"+strconv.Itoa(int(s.GetTimestamp()))+"}")
        modify = true
      }
    }

    // modify when is required
    if modify {
      newclusternode = clusternode.CopyClusterNode()
      newclusternode.SetStatus(s.GetStatus())
      newclusternode.AddService(s.CopyServiceObject())
      newclusternode.IncreaseTimestamp()
      cluster.AddNode(newclusternode)

      if err, clusterservice = NewClusterService(s.GetName()); err != nil {
        msg := "(ClusterEngine::handleServiceObject) "+err.Error()
        env.Output.WriteChDebug(msg)
        return errors.New(msg)
      }
      clusterservice.SetStatus(s.GetStatus())
      clusterservice.AddServiceNode(nodename,s.CopyServiceObject())
      cluster.AddService(clusterservice)

      env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Cluster update {name:'"+s.GetName()+"', status:"+strconv.Itoa(s.GetStatus())+", timestamp:"+strconv.Itoa(int(s.GetTimestamp()))+"}")
      if err := c.updateCluster(env.Whoami(),cluster); err != nil {
        env.Output.WriteChError("(ClusterEngine::handleServiceObject) "+err.Error())
      }
    } else {
      env.Output.WriteChError("(ClusterEngine::handleServiceObject) Service '"+s.GetName()+"' will not be analized to update cluster")
    }
  }
  return nil
}
//
//# handleServiceObject: handles the incoming serviceObject messages
func (c *ClusterEngine) handleClusterMessage(data []byte) error {
  env.Output.WriteChDebug("(ClusterEngine::handleClusterMessage)")
  var messageData *Cluster
  //
  // Decoding received []byte to a cluster message
  if err, message := DecodeClusterMessage(data); err != nil {
    // When the data could not be decoded an error is thrown
    env.Output.WriteChError("(ClusterEngine::handleClusterMessage) "+err.Error())
    return errors.New("(ClusterEngine::handleClusterMessage) "+err.Error())
  } else {
    env.Output.WriteChDebug("(ClusterEngine::handleClusterMessage) Message arrived {from_node:'"+message.GetFrom()+"', from_timestamp:"+strconv.Itoa(int(message.GetTimestamp()))+"}")
    cluster := c.GetCluster()

    if err, messageData = DecodeData(message.GetData()); err != nil {
      msg := "(ClusterEngine::handleClusterMessage) "+err.Error()
      env.Output.WriteChError(msg)
      return errors.New(msg)
    }
    //
    // Read data from the message: 
    // get the reviced message timestamp
    messageTimestamp := message.GetTimestamp()
    // get sender clusternode's stored data
    if err, clusternode := cluster.GetNode(message.GetFrom()); err != nil {
      env.Output.WriteChDebug("(ClusterEngine::handleClusterMessage) Received a message from new node '"+message.GetFrom()+"', and will be added on cluster")
      //
      // update information from nodes
      if err := c.updateCluster(message.GetFrom(),messageData); err != nil {
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
        env.Output.WriteChDebug("(ClusterEngine::handleServiceObject) Before updateCluster")
        if err := c.updateCluster(clusternode.GetName(),messageData); err != nil {
          env.Output.WriteChError("(ClusterEngine::handleClusterMessage) "+err.Error())
        }
      }
    }
    
  }
  return nil
}
//
//# updateCluster: method to update stored nodes information
func (c *ClusterEngine) updateCluster(from string, clusterFrom *Cluster) error {
  env.Output.WriteChDebug("(ClusterEngine::updateCluster)")
  modify := false
  cluster := c.GetCluster()

  for _,nodeFrom := range clusterFrom.GetNodes(){
    //
    // get stored node's data from cluster
    env.Output.WriteChDebug("(ClusterEngine::updateCluster) Received data {node:'"+nodeFrom.GetName()+"', node_status:"+strconv.Itoa(nodeFrom.GetStatus())+", node_timestamp:"+strconv.Itoa(int(nodeFrom.GetTimestamp()))+"}")
    // node has not been already stored
    if err, nodeStored := cluster.GetNode(nodeFrom.GetName()); err != nil {
      env.Output.WriteChDebug("(ClusterEngine::updateCluster) Node will be added to cluster {node:'"+nodeFrom.GetName()+"', node_status:"+strconv.Itoa(nodeFrom.GetStatus())+", node_timestamp:"+strconv.Itoa(int(nodeFrom.GetTimestamp()))+"}")
      modify = true
    // node already exist on cluster
    } else {
      // update stored node when received timestamp from node is greater
      env.Output.WriteChDebug("(ClusterEngine::updateCluster) Node '"+nodeFrom.GetName()+"' already exist on cluster")
      if err, nodeFromService := nodeFrom.HasService(nodeFrom.GetName()); err != nil {
        env.Output.WriteChWarn("(ClusterEngine::updateCluster) Could not load service from '"+nodeFrom.GetName()+"'."+err.Error())
      } else {
        // get number of checks from the service. Is possible to detect a restart when number of checks is lower than timestamp
        numberOfChecks := len(nodeFromService.GetChecks())
        //env.Output.WriteChDebug("(ClusterEngine::updateCluster) '"+nodeFrom.GetName()+"' numberOfChecks:"+strconv.Itoa(numberOfChecks))
        //env.Output.WriteChDebug("(ClusterEngine::updateCluster) '"+nodeFrom.GetName()+"' nodeFrom timestamp:"+strconv.Itoa(int(nodeFrom.GetTimestamp())))
        //env.Output.WriteChDebug("(ClusterEngine::updateCluster) '"+nodeStored.GetName()+"' nodeStored timestamp:"+strconv.Itoa(int(nodeStored.GetTimestamp())))

        // validate if modification is required
        if ( nodeFrom.GetTimestamp() > nodeStored.GetTimestamp()) || int(nodeFrom.GetTimestamp()) <= numberOfChecks {
          env.Output.WriteChDebug("(ClusterEngine::updateCluster) Node will be updated to cluster {node:'"+nodeFrom.GetName()+"', node_status:"+strconv.Itoa(nodeFrom.GetStatus())+", node_timestamp:"+strconv.Itoa(int(nodeFrom.GetTimestamp()))+"}")
          modify = true
        } else {
          env.Output.WriteChDebug("(ClusterEngine::updateCluster) Node won't be updated {node:'"+nodeFrom.GetName()+"'}")
        }
      }
    }
    //
    // update services related with node when node has been updated
    // it could be done for each node on cluster
    if modify {
      env.Output.WriteChDebug("(ClusterEngine::updateCluster) Modifing data on node {node:'"+nodeFrom.GetName()+"', node_status:"+strconv.Itoa(nodeFrom.GetStatus())+", node_timestamp:"+strconv.Itoa(int(nodeFrom.GetTimestamp()))+"}")
      // add the node to cluster
      cluster.AddNode(nodeFrom)
      if nodeFrom.GetCandidateForDeletion() {
        if err, deletable := cluster.ConsensusForDeletion(nodeFrom.GetName(),from,false); err != nil {
          env.Output.WriteChError("(ClusterEngine::StartClusterSync) "+err.Error())
        } else {
          if deletable {
            //TODO
          }
        }
      }
      //
      // for each service defined on node, check it's status
      for _,serviceFrom := range nodeFrom.GetNodeServices(){
        env.Output.WriteChDebug("(ClusterEngine::updateCluster) Review service {node:'"+nodeFrom.GetName()+"', service:'"+serviceFrom.GetName()+"', service_status:"+strconv.Itoa(serviceFrom.GetStatus())+", service_timestamp:"+strconv.Itoa(int(serviceFrom.GetTimestamp()))+"}")
        // get cluster service stored for serviceFrom
        if err, clusterService := cluster.GetService(serviceFrom.GetName()); err != nil {
          if err, clusterService = NewClusterService(serviceFrom.GetName()); err != nil {
            env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
          } else {
            env.Output.WriteChDebug("(ClusterEngine::updateCluster) Add service {node:'"+nodeFrom.GetName()+"', service:'"+serviceFrom.GetName()+"', service_status:"+strconv.Itoa(serviceFrom.GetStatus())+", service_timestamp:"+strconv.Itoa(int(serviceFrom.GetTimestamp()))+"}")
            // add service on cluster
            c.AddService(clusterService)
            // relate service to node
            clusterService.AddServiceNode(nodeFrom.GetName(),serviceFrom)
          }
        } else {
          // get stored data for the cluster service
          if err, serviceStored := clusterService.GetServiceNode(nodeFrom.GetName()); err != nil{
            env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
          } else {
            //env.Output.WriteChDebug("(ClusterEngine::updateCluster) '"+serviceFrom.GetName()+"' serviceFrom timestamp:"+strconv.Itoa(int(serviceFrom.GetTimestamp())))
            //env.Output.WriteChDebug("(ClusterEngine::updateCluster) '"+serviceFrom.GetName()+"' serviceStored timestamp:"+strconv.Itoa(int(serviceStored.GetTimestamp())))
                  
            // update service whem serviceFrom timestamp is greater than the stored one
            if serviceFrom.GetTimestamp() > serviceStored.GetTimestamp() {
              env.Output.WriteChDebug("(ClusterEngine::updateCluster) Update service {node:'"+nodeFrom.GetName()+"', service:'"+serviceFrom.GetName()+"', service_status:"+strconv.Itoa(serviceFrom.GetStatus())+", service_timestamp:"+strconv.Itoa(int(serviceFrom.GetTimestamp()))+"}")
              // relate node
              clusterService.AddServiceNode(nodeFrom.GetName(),serviceFrom)
            }
          }
        }
      }
    }
  }

  //
  // On any change, increase timestamp and deploy clusterMessage to cluster
  // it could be done once for cluster update
  if modify {
    env.Output.WriteChDebug("(ClusterEngine::updateCluster) Cluster status has changed")
    if err, clusternode := cluster.GetNode(env.Whoami()); err != nil {
      env.Output.WriteChError("(ClusterEngine::updateCluster) Cluster status has changed but ClusterMessage could not be deploy due an error: "+err.Error())
    } else {
      // Prepare the message to be sent to the cluster
      if err, message := NewClusterMessage(env.Setup.Hostname, clusternode.GetTimestamp(),cluster); err == nil {
        env.Output.WriteChDebug("(ClusterEngine::updateCluster) Cluster message ready to be sent {from_node:"+clusternode.GetName()+", from_timestamp:"+strconv.Itoa(int(clusternode.GetTimestamp()))+"}")
        // transform data as []byte before send it
        if err, messageBytes := utils.InterfaceToBytes(message); err == nil {
          env.Output.WriteChDebug("(ClusterEngine::updateCluster) Send data to outputChannels")
          // send data to the outputChannels
          c.SendData(messageBytes)
        } else {
          env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
        }
      }

      // if err, messageData = utils.ObjectToJsonByte(cluster); err != nil {
      //   msg := "(ClusterEngine::updateCluster) "+err.Error()
      //   env.Output.WriteChError(msg)
      //   return errors.New(msg)
      // }
      // if err, message := NewClusterMessage(env.Whoami(), clusternode.GetTimestamp(),messageData); err == nil {
      //   env.Output.WriteChDebug("(ClusterEngine::updateCluster) New message ready to be sent")
      //   // transform data as []byte before send it
      //   if err, messageBytes := utils.ObjectToJsonByte(message); err == nil {
      //     env.Output.WriteChDebug("(ClusterEngine::updateCluster) Send data to outputChannels")
      //     // send data to the outputChannels
      //     c.SendData(messageBytes)
      //   } else {
      //     env.Output.WriteChError("(ClusterEngine::updateCluster) "+err.Error())
      //   }    
      // }
    }
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
  
  return utils.ObjectToJsonByte(c.GetCluster())
}
//
//# GetClusterNodes: get all nodes from cluster
func (c *ClusterEngine) GetClusterNodes() (error,[]byte) {
  env.Output.WriteChDebug("(ClusterEngine::GetClusterNodes)")
  var cluster *Cluster

  if cluster = c.GetCluster(); cluster == nil {
    return errors.New("(ClusterEngine::GetClusterNodes) Cluster object has not been initialized"),nil
  }

  return utils.ObjectToJsonByte(cluster.GetNodes())
}
//
//# GetClusterServices: get all nodes from cluster
func (c *ClusterEngine) GetClusterServices() (error,[]byte) {
  env.Output.WriteChDebug("(ClusterEngine::GetClusterServices)")
  var cluster *Cluster

  if cluster = c.GetCluster(); cluster == nil {
    return errors.New("(ClusterEngine::GetClusterServices) Cluster object has not been initialized"),nil
  }

  return utils.ObjectToJsonByte(cluster.GetServices())
}

//
// Sync
//-----------------------------------------------------------------------

//# SelectNodesForSync
func (c *ClusterEngine) SelectNodesForSync() (error, []string) {
  env.Output.WriteChDebug("(ClusterEngine::SelectNodesForSync)")
  cluster := c.GetCluster()
  
  //selectedNodes := make(map[string]*ClusterNode)
  numberOfNodes := len(cluster.GetNodes())

  // current nodes is alone
    // if numberOfNodes <= 1 {
    //   env.Output.WriteChDebug("(ClusterEngine::SelectNodesForSync) No nodes to sync")   
    //   return nil, nil
    // }
  // number of nodes to sync
  numberOfSelectedNodes := int(math.Log2(float64(numberOfNodes)))
  env.Output.WriteChDebug("(ClusterEngine::SelectNodesForSync) "+strconv.Itoa(numberOfSelectedNodes))
  // selected nodes, current node won't be sync
  clusternodes := make([]string,numberOfNodes-1)
  whoami := env.Whoami()

  rand.Seed(time.Now().UTC().UnixNano())
  it := 0
  for nodename,_ := range cluster.GetNodes() {
    if nodename != whoami {
      env.Output.WriteChDebug("(ClusterEngine::SelectNodesForSync) Selected node {node:'"+nodename+"'}")
      clusternodes[it] = nodename
      it++
    }
  }
  return nil, clusternodes
}

//
//# StartClusterSync
func (c *ClusterEngine) StartClusterSync() error {
  env.Output.WriteChDebug("(ClusterEngine::StartClusterSync) Starting cluster Sync")
  cluster := c.GetCluster()
  
  go func() {
    defer close (c.inputSyncChannel)
    for{
      select{
      case message := <-c.inputSyncChannel:
        env.Output.WriteChDebug("(ClusterEngine::StartClusterSync) New message to Sync")
        if err, nodes := c.SelectNodesForSync(); err == nil {
          for _,nodename := range nodes {
            env.Output.WriteChDebug("(ClusterEngine::StartClusterSync) Selected node "+nodename)
            go func() {
              if err, clusternode := cluster.GetNode(nodename); err == nil {
                env.Output.WriteChError("(ClusterEngine::StartClusterSync) Selected node to send sync message {from_node:"+env.Whoami()+"', to_node:'"+clusternode.GetName()+"'}")
                e := strings.Split(clusternode.GetURL(),"://")
                env.Output.WriteChDebug("(ClusterEngine::StartClusterSync) endpoint: "+e[1])
                if err := utils.CheckEndpoint("tcp",e[1]); err != nil {
                  env.Output.WriteChError("(ClusterEngine::StartClusterSync) "+err.Error())
                  clusternode.SetCandidateForDeletion(true)
                  if err, deletable := cluster.ConsensusForDeletion(nodename,env.Whoami(),false); err != nil {
                    env.Output.WriteChError("(ClusterEngine::StartClusterSync) "+err.Error())
                  } else {
                    if deletable {
                      //TODO
                    }
                  }
                } else {
                  if clusternode.GetCandidateForDeletion() == true {
                    clusternode.SetCandidateForDeletion(false)
                  }
                  env.Output.WriteChDebug("(ClusterEngine::StartClusterSync) Send sync message {from_node:"+env.Whoami()+"', to_node:'"+clusternode.GetName()+"'}")
                  if err := c.SendSyncMessage(clusternode.GetURL(),message); err != nil {
                    env.Output.WriteChError("(ClusterEngine::StartClusterSync) "+err.Error())
                  }
                }
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
//# SendSyncMessage
func (c *ClusterEngine) SendSyncMessage(url string, message []byte) error {
  
  uri := "/api/cluster"
  url = url + uri
  env.Output.WriteChDebug("(ClusterEngine::SendSyncMessage) Message to '"+url+"'")

  req, err := http.NewRequest("PUT", url, bytes.NewBuffer(message))
  req.Header.Set("X-Verdmell-From", env.Whoami())
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)


  if err != nil {
    msg := "(ClusterEngine::SendSyncMessage) "+err.Error()
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
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

import(
	"errors"
	"verdmell/service"
	"verdmell/utils"
)

//#
//#
//# ClusterMessage struct:
//# ClusterMessage is defined by a node name and its URL
type ClusterMessage struct{
	From string `json:"from"`
	Timestamp int64 `json:"timestamp"`
	Data interface{}	`json:"message"`
}

//
//# NewClusterMessage: return a CheckMessage
func NewClusterMessage(f string, t int64, i interface{}) (error, *ClusterMessage){
	env.Output.WriteChDebug("(ClusterMessage::NewClusterMessage) Creating a new ClusterMessage")

	// check for a valid timestamp
	if t < 0 {
		msg := "(ClusterMessage::NewClusterMessage) Not valid timestamp"
		env.Output.WriteChError(msg)
		return errors.New(msg),nil
	}
	switch i.(type) {
	case *service.ServiceObject:
		env.Output.WriteChDebug("(ClusterMessage::NewClusterMessage) New message for ServiceObject")
	case *ClusterNode:
		env.Output.WriteChDebug("(ClusterMessage::NewClusterMessage) New message for ServiceObject")
	default:
		msg := "(ClusterMessage::NewClusterMessage) Not valid data type for message"
		env.Output.WriteChError(msg)
		return errors.New(msg),nil
	}
	
	message := &ClusterMessage{
		From: f,
		Timestamp: t,
		Data: i,
	}

	return nil,	message
}

//
//# SetFrom: attribute from ClusterMessage
func (m *ClusterMessage) SetFrom(f string) {
	env.Output.WriteChDebug("(ClusterMessage::SetFrom)")
	m.From = f
}
//
//# SetTimestamp: attribute from ClusterMessage
func (m *ClusterMessage) SetTimestamp(t int64) {
	env.Output.WriteChDebug("(ClusterMessage::SetTimestamp)")	
	m.Timestamp = t
}
//
//# SetData: attribute from ClusterMessage
func (m *ClusterMessage) SetData(d interface{}) {
	env.Output.WriteChDebug("(ClusterMessage::SetData)")
	m.Data = d
}

//
//# GetFrom: attribute from ClusterMessage
func (m *ClusterMessage) GetFrom() string {
	env.Output.WriteChDebug("(ClusterMessage::GetFrom)")
	return m.From
}
//
//# GetTimestamp: attribute from ClusterMessage
func (m *ClusterMessage) GetTimestamp() int64 {
	env.Output.WriteChDebug("(ClusterMessage::GetTimestamp)")	
	return m.Timestamp
}
//
//# GetData: attribute from ClusterMessage
func (m *ClusterMessage) GetData() interface{} {
	env.Output.WriteChDebug("(ClusterMessage::GetData)")
	return m.Data
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# GetData: attribute from ClusterMessage
func (m *ClusterMessage) EncodeData() (error, interface{}) {
	env.Output.WriteChDebug("(ClusterMessage::DecodeData)")
	return utils.InterfaceToBytes(m.GetData())
}

//#
//# Common methods
//#---------------------------------------------------------------------

//
//# String: convert a ClusterMessage object to string
func (m *ClusterMessage) String() string {
  if err, str := utils.ObjectToJsonString(m); err != nil{
    return err.Error()
  } else{
    return str
  }
}

//#######################################################################################################
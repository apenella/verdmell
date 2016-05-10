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
	"bytes"
	"encoding/gob"
	"verdmell/utils"
)

//#
//#
//# ClusterMessage struct:
//# ClusterMessage is defined by a node name and its URL
type ClusterMessage struct{
	From string `json:"from"`
	Timestamp int64 `json:"timestamp"`
	Data []byte `json:"data"`
}

//
//# NewClusterMessage: return a CheckMessage
func NewClusterMessage(f string, t int64, i interface{}) (error, *ClusterMessage){
	env.Output.WriteChDebug("(ClusterMessage::NewClusterMessage) Creating a new ClusterMessage")

	message := new(ClusterMessage)

	// check for a valid timestamp
	if t < 0 {
		msg := "(ClusterMessage::NewClusterMessage) Not valid timestamp"
		env.Output.WriteChError(msg)
		return errors.New(msg),nil
	}
	// check the cluster state
	if i == nil {
		msg := "(ClusterMessage::NewClusterMessage) The message requieres not empty Cluster"
		env.Output.WriteChError(msg)
		return errors.New(msg),nil	
	}

	switch i.(type) {
		case []byte:
			env.Output.WriteChDebug("(ClusterMessage::NewClusterMessage) []byte type arrived")
			message.SetData(i.([]byte))
			break;
		case *Cluster:
			env.Output.WriteChDebug("(ClusterMessage::NewClusterMessage) *Cluster type arrived")
			if err, data := utils.InterfaceToBytes(i); err != nil {
				msg := "(ClusterMessage::NewClusterMessage) "+err.Error()
				env.Output.WriteChError(msg)
				return errors.New(msg),nil	
			} else {
				message.SetData(data)
			}
			break;
		default:
			msg := "(ClusterMessage::NewClusterMessage) Not valid type to construct a message"
			env.Output.WriteChError(msg)
			return errors.New(msg),nil
	}

	message.SetFrom(f)
	message.SetTimestamp(t)

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
func (m *ClusterMessage) SetData(d []byte) {
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
func (m *ClusterMessage) GetData() []byte {
	env.Output.WriteChDebug("(ClusterMessage::GetData)")
	return m.Data
}

//#
//# Specific methods
//#---------------------------------------------------------------------

//
//# EncodeData: attribute from ClusterMessage
func EncodeData(m *Cluster) (error, []byte) {
	env.Output.WriteChDebug("(ClusterMessage::EncodeData)")
	return utils.InterfaceToBytes(m)
}
//
//# DecodeData: return a byted data
func DecodeData(data []byte) (error, *Cluster) {
	env.Output.WriteChDebug("(ClusterMessage::DecodeData)")
	message := new(Cluster)
 
 	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	if err := dec.Decode(message); err != nil {
		return errors.New("(ClusterMessage::DecodeData) "+err.Error()),nil
	}

	return nil, message
}
//
//# DecodeClusterMessage: return a *ClusterMessaged data
func DecodeClusterMessage(data []byte) (error, *ClusterMessage) {
	env.Output.WriteChDebug("(ClusterMessage::DecodeClusterMessage)")
	message := new(ClusterMessage)
 
 	buffer := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buffer)
	if err := dec.Decode(message); err != nil {
		return errors.New("(ClusterMessage::DecodeClusterMessage) "+err.Error()),nil
	}

	return nil, message
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
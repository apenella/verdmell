package utils

import (
  "encoding/gob"
  "bytes"
  "github.com/apenella/messageOutput"
)

//
//#InterfaceToBytes: convert a interface{} to a []byte
func InterfaceToBytes(key interface{}) (error, []byte) {
    var buf bytes.Buffer
    
    message.WriteDebug("(InterfaceToBytes)")
    
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(key)
    if err != nil {
        return err, nil
    }
    return nil, buf.Bytes()
}


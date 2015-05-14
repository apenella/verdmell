package utils

import (
  "encoding/gob"
  "bytes"
  "fmt"
  "github.com/apenella/messageOutput"
)

func GetInterfaceBytes(key interface{}) ([]byte, error) {
    var buf bytes.Buffer
    
    message.WriteDebug("(GetInterfaceBytes)")
    fmt.Println(key)
    
    enc := gob.NewEncoder(&buf)
    err := enc.Encode(key)
    if err != nil {
        return nil, err
    }
    return buf.Bytes(), nil
}

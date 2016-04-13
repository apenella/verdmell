package utils

import (
  "os"
  "encoding/json"
  "io/ioutil"
  "github.com/apenella/messageOutput"
)

//# 
//# functions for JSON
//#--------------------------------------------------------------------

//
//# LoadJSONFile: function to dump data from the file f to object
func LoadJSONFile(f string, object interface{}) error {
  file, e := ioutil.ReadFile(f)

  if e != nil {
    message.WriteError("(loadJSONFile) File error")
    os.Exit(1)
  }
  return json.Unmarshal(file, object)
}

//
//# ObjectToJsonString: converst any object to a json string
func ObjectToJsonString(object interface{}) (error, string) {
  if jsoned, err := json.Marshal(object); err != nil{
    return err,err.Error() 
  }else{
    return nil,string(jsoned)
  }
  //jsoned, _ := json.Marshal(object)
  //return string(jsoned)
}
//
//# ObjectToJsonString: converst any object to a json string
func ObjectToJsonStringPretty(object interface{}) (error,string) {
  if jsoned, err := json.MarshalIndent(object,"","	"); err != nil{
    return err,err.Error() 
  }else{
    return nil,string(jsoned)
  }
}
//
//# ObjectToJsonByte: converst any object to a json byte
func ObjectToJsonByte(object interface{}) (error, []byte) {
  jsoned, _ := json.Marshal(object)
  return nil, jsoned
}

//#######################################################################################################
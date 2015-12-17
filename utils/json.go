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
func LoadJSONFile(f string, object interface{}) {
  file, e := ioutil.ReadFile(f)

  if e != nil {
    message.WriteError("(loadJSONFile) File error")
    os.Exit(1)
  }
  json.Unmarshal(file, object)
}

//
//# ObjectToJsonString: converst any object to a json string
func ObjectToJsonString(object interface{}) string {
  jsoned, _ := json.Marshal(object)
  return string(jsoned)
}
//
//# ObjectToJsonString: converst any object to a json string
func ObjectToJsonStringPretty(object interface{}) string {
  jsoned, _ := json.MarshalIndent(object,"","	")
  return string(jsoned)
}
//
//# ObjectToJsonByte: converst any object to a json byte
func ObjectToJsonByte(object interface{}) []byte {
  jsoned, _ := json.Marshal(object)
  return jsoned
}

//#######################################################################################################
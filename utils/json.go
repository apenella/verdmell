package utils

import (
  "errors"
  "encoding/json"
  "io/ioutil"
)

//# 
//# functions for JSON
//#--------------------------------------------------------------------

//
//# LoadJSONFile: function to dump data from the file f to object
func LoadJSONFile(f string, object interface{}) error {
  
  if file, err := ioutil.ReadFile(f); err != nil {
    return errors.New("(loadJSONFile) Error on loading file '"+f+"'. "+err.Error())
  } else {
    if err:= json.Unmarshal(file, object); err != nil{
      return errors.New("(loadJSONFile) Error on loading file '"+f+"'. "+err.Error())
    }
  }
  return nil
}

//
//# ObjectToJsonString: converst any object to a json string
func ObjectToJsonString(object interface{}) (error, string) {
  if jsoned, err := json.Marshal(object); err != nil{
    return err, err.Error() 
  }else{
    return nil,string(jsoned)
  }
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
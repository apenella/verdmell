package utils

import (
  "os"
  "encoding/json"
  "io/ioutil"
  "github.com/apenella/messageOutput"
)

// functions for JSON
// function to dump data from the file f to object
func LoadJSONFile(f string, object interface{}) {
    
  file, e := ioutil.ReadFile(f)

  if e != nil {
    message.WriteError("(loadJSONFile) File error")
    os.Exit(1)
  }

  json.Unmarshal(file, object)

}

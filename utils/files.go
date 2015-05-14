package utils

import (
  "os"
  "io/ioutil"
  "github.com/apenella/messageOutput"
)

// OS functions

func FileExist(file string) error {
	_, err := os.Stat(file);
	return err
}

// get all files from folder
func GetFolderFiles(folder string) []os.FileInfo {
  files,err := ioutil.ReadDir(folder)

  if err != nil {
    message.WriteError("(getFolderFiles) Read dir error")
    os.Exit(1)
  } 

  return files
}

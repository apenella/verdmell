package agent

import (
    "fmt"
    "strings"
)

// Create a new type for a list of Strings
type StringList []string

// Implement the flag.Value interface
func (s *StringList) String() string {
    return fmt.Sprintf("%v", *s)
}
func (s *StringList) Set(value string) error {
    *s = strings.Split(value, ",")
    return nil
}
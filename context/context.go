/*
Package context contains execution data details
*/
package context

import (
	"io"
	"verdmell/utils"

	"github.com/apenella/messageOutput"
)

//
// Context
type Context struct {
	// environment
	UI io.Writer `json:"-"`
	// Logger output manager
	Logger *message.Message `json: "-"`
	// host to anchor to server mode
	Host string `json: "host"`
	// port to anchor to server mode
	Port int `json: "port"`
	// nodes that belongs to cluster
	Cluster []string `json: "cluster"`
	// Loglevel for agent
	Loglevel int `json: "loglevel"`
	// Configfile is configuration file name
	Configfile string `json: "configuration_file"`
	// Configdir directory to place configuration
	Configdir string `json: "configuration_dir"`
}

// String method transform the Configuration to string
func (ctx *Context) String() string {
	var err error
	var str string

	str, err = utils.ObjectToJSONString(ctx)
	if err != nil {
		return err.Error()
	}

	return str
}

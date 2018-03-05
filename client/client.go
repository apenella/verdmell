package client

import (
	"github.com/mitchellh/cli"
)

// ClientWorker is an interface which represents the element that does the work
type ClientWorker interface {
	Run() error
}

// Client the struct represents an engine which is responsable of the interaction between user and daemon.
// Client works is done by an implementation of clientWorker interface. Each kind of user interaction is implemented by a diferent client worker.
type Client struct {
	ID uint `json: "id"`
	Worker ClientWorker
	Ui cli.Ui
}

// Init
func (c *Client) Init() error {
	return nil
}


// Run is responsable to run worker
func (c *Client) Run() error {
	err := c.Worker.Run()
	if err != nil {
		return err
	}

	return nil
}

// Stop
func (c *Client) Stop() error {
	return nil
}

// Subscribe
func (c *Client) Subscribe(o chan interface{}, desc string) error {
	return nil
}

// Status
func (c *Client) Status() error {
	return nil
}

// SayHi
func (c *Client) SayHi() {}

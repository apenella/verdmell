package client


/*
	ClientWorker defines those objects which could work as a client 
*/
type ClientWorker interface {
	Run() error
}

/*
	Client the struct for client engine
*/
type Client struct {
	ID uint `json: "id"`
	Worker ClientWorker
}

/*Init*/
func (c *Client) Init() error {
	return nil
}

/*
	Run is responsable to run worker
*/
func (c *Client) Run() error {
	err := c.Worker.Run()
	if err != nil {
		return err
	}

	return nil
}

/*Stop*/
func (c *Client) Stop() error {
	return nil
}

/*Subscribe*/
func (c *Client) Subscribe(o chan interface{}, desc string) error {
	return nil
}

/*Status*/
func (c *Client) Status() error {
	return nil
}

/*SayHi*/
func (c *Client) SayHi() {}

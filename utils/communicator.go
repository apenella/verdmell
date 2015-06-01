package utils

type Communicator interface{
	Send()
	Receive()
}

type Send interface{
	Out() <-chan interface{}
	Close() 
}
type Receive interface{
	In() chan <-interface{}
}

func WrapCommunicator(ch interface{}) {
	
}
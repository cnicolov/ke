package connection

import (
	"io"

	"kego.io/editor/shared/messages"
	"kego.io/ke"
	"kego.io/kerr"
	"kego.io/system"
)

// Send sends a message to the connection.
func (c *Conn) Send(message messages.MessageInterface) {
	c.outwg.Add(1)
	c.out <- message
}

// Request sends a message and expects a response
func (c *Conn) Request(message messages.MessageInterface) chan messages.MessageInterface {
	responseChannel := c.sendRequestAndReturnResponseChannel(message)
	outputChannel := c.applyTimeout(TIMEOUT, responseChannel)
	return outputChannel
}
func (c *Conn) sendRequestAndReturnResponseChannel(message messages.MessageInterface) chan messages.MessageInterface {
	responseChannel := make(chan messages.MessageInterface)
	c.requests[message.GetMessage().Guid.Value] = responseChannel
	c.outwg.Add(1)
	c.out <- message
	return responseChannel
}

// Respond sends a message as a reply to a request
func (c *Conn) Respond(message messages.MessageInterface, requestGuid string) {
	message.GetMessage().Request = system.NewString(requestGuid)
	c.outwg.Add(1)
	c.out <- message
}

func (c *Conn) sender() error {
	for {
		if err := c.senderInternal(); err != nil {
			return kerr.New("RCWDBROWAN", err, "senderInternal")
		}
	}
}
func (c *Conn) senderInternal() error {
	m := <-c.out
	defer c.outwg.Done()
	if err := ke.NewEncoder(c.socket).Encode(m); err != nil {
		if err == io.EOF {
			// Closing the fail channel exits the app gracefully
			if c.debug {
				c.fail <- kerr.New("HHLJYVNLJM", nil, "Connection closed")
				return nil
			}
			close(c.fail)
			return nil
		}
		return kerr.New("WIUXNWRXCQ", err, "Encode")
	}
	return nil
}

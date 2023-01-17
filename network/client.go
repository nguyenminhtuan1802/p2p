package network

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"tnguyen/blockchainexample/utility"
)

type Client struct {
	Protocol   string
	Address    string
	Connection *net.Conn
}

func ClientConstructor(protocol string, address string) *Client {
	client := new(Client)
	client.Protocol = protocol
	client.Address = address
	return client
}

func (c *Client) Connect() {
	connection, err := net.Dial(c.Protocol, c.Address)
	utility.CheckError(err, "[CLIENT] Connect Exception")

	c.Connection = &connection
	//fmt.Println("[CLIENT] Connect Successful")
}

func (c *Client) Send(msg string) {
	// Wait for valid connection
	for (c.Connection != nil) {
		break
	}
	err := gob.NewEncoder(*c.Connection).Encode(msg)
	//fmt.Println("[CLIENT ", (*c.Connection).LocalAddr().String(), "] Sent: ", msg)
	utility.CheckError(err, "[CLIENT] Send Exception")
}

func (c *Client) ReceiveOnce() interface{} {
	// Wait for valid connection
	for (c.Connection != nil) {
		break
	}

	// receive 1 message and return
	var msg string

	for {
		err := gob.NewDecoder(*c.Connection).Decode(&msg)
		if err != nil {
			//fmt.Println("[SERVER] Receveid from client error:]", err)
		} else {			
			res, err := strconv.Atoi(msg)
			if err == nil {				
				//fmt.Println("[CLIENT ", (*c.Connection).LocalAddr().String(), "] Received:", msg)
				return res
			} else {
				return msg
			}
		}
	}
}

func (c *Client) Receive() {
	// Wait for valid connection
	for (c.Connection != nil) {
		break
	}
	var msg string
	for {
		err := gob.NewDecoder(*c.Connection).Decode(&msg)
		if err != nil {
			//fmt.Println("[SERVER] Receveid from client error:]", err)
		} else {
			fmt.Println("[CLIENT ", c.Address, "] Received:", msg)
		}
	}
}

func (c *Client) Disconnect() {
	err := (*c.Connection).Close()
	utility.CheckError(err, "[CLIENT] Disconnect Exception")
}

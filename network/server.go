package network

import (
	"encoding/gob"
	"fmt"
	"net"
	"strconv"
	"regexp"
	"tnguyen/blockchainexample/utility"
)

type Server struct {
	Protocol string
	Address  string
	Listener *net.Listener
	Connections []*net.Conn
	Channel	chan int
}

func ServerConstructor(protocol string, address string, ch chan int) *Server {
	server := new(Server)
	server.Protocol = protocol
	server.Address = address
	server.Channel = ch
	listener, err := net.Listen(server.Protocol, server.Address)
	utility.CheckError(err, "[SERVER] Start Exception")
	server.Listener = &listener
	return server
}

func (s *Server) Start() {
	for {
		if (s.Listener != nil) {
			break
		}
	}

	fmt.Println("[SERVER] Starting")
	for {
		// accept a connection
		connection, err := (*s.Listener).Accept()
		if err != nil {
			// Checked if server has shut down
			pattern := "use of closed network connection"
			str := err.Error()
			res, _ := regexp.MatchString(pattern, str)
			if res {
				return
			}
			fmt.Println("[SERVER] Accept Exception", err)			
			continue
		}
		// handle the connection
		go s.OnNewConnection(&connection)
	}
}

func (s *Server) OnNewConnection(c *net.Conn) {
	fmt.Println("[SERVER ", (*s.Listener).Addr().String(), "] Accept new connection from [CLIENT ", (*c).LocalAddr().String() ,"]:")
	s.Connections = append(s.Connections, c)
	go s.ReceiveFromClient(c)
}

func (s *Server) SendToClient(c *net.Conn, msg string) {
	// send the message
	err := gob.NewEncoder(*c).Encode(&msg)
	if err != nil {
		fmt.Println("[SERVER ", s.Address, "] Sent to client error:", err)
	} else {
		//fmt.Println("[SERVER ", (*s.Listener).Addr().String(), "] Sent to [CLIENT ", (*c).LocalAddr().String() ,"] :", msg)
	}
}

func (s *Server) ReceiveFromClient(c *net.Conn) {
	// receive the message
	var msg string
	
	for {
		err := gob.NewDecoder(*c).Decode(&msg)
		if err == nil {
			
			if (msg == "QUERY") {
				nodeData := <-s.Channel
				go s.SendToClient(c, strconv.Itoa(nodeData))
			} else {
				//fmt.Println("[SERVER ", (*s.Listener).Addr().String(), "] Received from [CLIENT ", (*c).LocalAddr().String(), " ] ", msg)
				go s.SendToClient(c, "Reply to" + msg)
			}
		}
	}
}

func (s *Server) DisconnectClient(c *net.Conn) {
	err := (*c).Close()
	utility.CheckError(err, "[SERVER] Disconnect Client Exception ")
}

func (s *Server) ShutDown() {
	fmt.Println("[SERVER] Shuting Down")
	for _,c := range s.Connections {
		s.DisconnectClient(c)
	}
	(*s.Listener).Close()
}

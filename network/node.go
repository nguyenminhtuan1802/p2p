package network

import (
	"fmt"
	"strconv"
	"tnguyen/blockchainexample/transaction"
)

// A node represents a peer in P2P network
// Node contains a server that listens to peer node request 
// Clients in a node are used to send request to this node's peers
// Note: Channel provides data to respond to peers' requests. In this
//       example the data is transaction id, but channel can support other data types
type Node struct {
	Protocol 				string
	Id       				int
	Server   				*Server
	Clients  				map[int]*Client
	PublicTransactionList	*transaction.TransactionList
	PrivateTransactionList	*transaction.TransactionList
	Channel					chan int
}

var addresses = map[int]string{
	1:  "localhost:60000",
	2:  "localhost:60001",
	3:  "localhost:60002",
	4:  "localhost:60003",
	5:  "localhost:60004",
	6:  "localhost:60005",
	7:  "localhost:60006",
	8:  "localhost:60007",
	9:  "localhost:60008",
	10: "localhost:60090",
	11: "localhost:60010",
	12: "localhost:60011",
	13: "localhost:60012",
	14: "localhost:60013",
	15: "localhost:60014",
	16: "localhost:60015",
	17: "localhost:60016",
	18: "localhost:60017",
	19: "localhost:60018",
	20: "localhost:60019",
}

func NodeConstructor(protocol string, id int, publicTrans *transaction.TransactionList, privateTrans *transaction.TransactionList) *Node {
	node := new(Node)
	node.Protocol = protocol
	node.Id = id
	node.Clients = make(map[int]*Client)
	node.Channel = make(chan int, 1)

	node.PrivateTransactionList = privateTrans
	node.PublicTransactionList = publicTrans

	return node
}

func (n *Node) Start() {	
	n.Server = ServerConstructor(n.Protocol, addresses[n.Id], n.Channel)
	go n.DataStream()
	go n.Server.Start()
}

// Stream data thru channel to respond to peers' request
func (n *Node) DataStream() {
	for {
		n.Channel <- n.PrivateTransactionList.GetCurrentTransaction().Id
	}
}

// Connect to another peer
func (n *Node) Connect(another *Node) {
	_, ok := n.Clients[another.Id]

	if ok {
		return
	}

	client := ClientConstructor(another.Protocol, addresses[another.Id])
	client.Connect()
	anotherClient := ClientConstructor(n.Protocol, addresses[n.Id])
	anotherClient.Connect()

	n.Clients[another.Id] = client
	another.Clients[n.Id] = anotherClient
}

func (n *Node) PingAllConnections() {
	var i = 0
	for _, v := range n.Clients {
		v.Send("Ping from Node " + strconv.Itoa(n.Id) + " to Node " + v.Address)
		i++
		go v.Receive()
	}
	fmt.Println("[NODE ", n.Id, "] has", i, " connections")
}

func (n *Node) ShutDown() {
	n.Server.ShutDown()
}
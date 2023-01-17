package network

import (
	"fmt"
	"strconv"
	"tnguyen/blockchainexample/transaction"
	"tnguyen/blockchainexample/utility"
	"time"
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
	Validator				ConsensusValidator
	FinishRunning			bool
}

func NodeConstructor(protocol string, id int,
	 publicTrans *transaction.TransactionList,
	 privateTrans *transaction.TransactionList,) *Node {
	node := new(Node)
	node.Protocol = protocol
	node.Id = id
	node.Clients = make(map[int]*Client)
	node.Channel = make(chan int, 1)

	node.PrivateTransactionList = privateTrans
	node.PublicTransactionList = publicTrans
	node.FinishRunning = false

	return node
}

func (n *Node) Start() {	
	n.Server = ServerConstructor(n.Protocol, utility.Addresses[n.Id], n)
	go n.Server.Start()
}

func (n *Node) Run() {
	time.Sleep(5 * time.Second)
	for {
		if n.Validator.Update() == false {
			fmt.Println("[NODE ", n.Id, "] DONE")
			for _,v := range n.PublicTransactionList.Transactions {
				fmt.Println("[NODE ", n.Id, "] :", v.Data)
			}
			n.FinishRunning = true
			break
		}
	}
}

// Connect to another peer
func (n *Node) ConnectToNode(another *Node) {
	_, ok := n.Clients[another.Id]

	if ok {
		return
	}

	client := ClientConstructor(another.Protocol, utility.Addresses[another.Id])
	client.Connect()
	anotherClient := ClientConstructor(n.Protocol, utility.Addresses[n.Id])
	anotherClient.Connect()

	n.Clients[another.Id] = client
	another.Clients[n.Id] = anotherClient
}

// Connect to another peer by id/address
func (n *Node) ConnectById(protocol string, id int) {
	_, ok := n.Clients[id]

	if ok {
		return
	}

	client := ClientConstructor(protocol, utility.Addresses[id])
	client.Connect()
	n.Clients[id] = client
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

func (n *Node) GetNextUndecidedTransaction() (int, *transaction.Transaction) {
	index := n.PrivateTransactionList.CurrentIdx
	tran := n.PrivateTransactionList.Transactions[index]
	return index, tran
}

func (n *Node) GetTransactionData(index int) int {
	publicTransaction := n.PublicTransactionList.GetTransactionData(index)
	if (publicTransaction != nil) {
		return publicTransaction.Data
	} else {
		privateTransaction := n.PrivateTransactionList.GetTransactionData(index)
		if (privateTransaction != nil) {
			return privateTransaction.Data
		} else {
			panic("Invalid index")
		}
	}
}

func (n *Node) HasFinished() bool {
	return n.FinishRunning
}

type ConsensusValidator interface {
	Update() bool
}
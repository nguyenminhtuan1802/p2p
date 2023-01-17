package test

import (
	//"fmt"
	"tnguyen/blockchainexample/transaction"
	"tnguyen/blockchainexample/consensus"
	"tnguyen/blockchainexample/network"
	"time"
)

type P2P struct {
	Nodes []*network.Node
}

func CreateP2PNetwork(size int, transactionLen int, key int, options *consensus.SnowballConsensusOptions) *P2P {
	net := new(P2P)
	net.Nodes = make([]*network.Node, 0)
	
	for i := 1; i <= size; i++ {
		len := transactionLen
		// create transaction list
		privateList := new(transaction.TransactionList)
		if i <= size/2 {
			for (len > 0) {
				privateList.Insert(key)
				len--
			}
		} else if i <= size * 3/4 {
			for (len > 0) {
				privateList.Insert(key + 1)
				len--
			}
		} else {
			for (len > 0) {
				privateList.Insert(i)
				len--
			}
		}

		publicList := new(transaction.TransactionList)

		node := network.NodeConstructor("tcp", i, publicList, privateList)
		validator := consensus.CreateSnowballConsensus(node, *options)
		node.Validator = validator
		net.Nodes = append(net.Nodes, node)
	}

	for _, n := range net.Nodes {
		time.Sleep(100 * time.Millisecond)
		n.Start()
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if j == i {
				continue
			}
			net.Nodes[i].Connect(net.Nodes[j])
		}
	}

	return net
}

func (net *P2P) Run() {
	for _, n := range net.Nodes {
		go n.Run()
	}
}

func (net *P2P) ShutDown() {
	for _,n := range net.Nodes {
		go n.ShutDown()
	}

	for _,n := range net.Nodes {
		for {
			if (n.Server.Running == false) {
				break
			}
		}
	}
}
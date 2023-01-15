package network

import (
	"fmt"
	"tnguyen/blockchainexample/transaction"
)

type P2P struct {
	Nodes []*Node
}

func CreateP2PNetwork(size int, transactionLen int, key int) *P2P {
	net := new(P2P)
	net.Nodes = make([]*Node, 0)
	
	for i := 1; i <= size; i++ {
		len := transactionLen
		// create transaction list
		privateList := new(transaction.TransactionList)
		if i%3 == 0 {
			for (len > 0) {
				privateList.Insert(key)
				len--
			}
		} else {
			for (len > 0) {
				privateList.Insert(i)
				len--
			}
		}

		publicList := new(transaction.TransactionList)

		node := NodeConstructor("tcp", i, publicList, privateList)
		net.Nodes = append(net.Nodes, node)
	}

	for _, n := range net.Nodes {
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

func (net *P2P) GetPublicTransactionList() []int {
	ans := make([]int, 0)
	node := net.Nodes[0]
	cur := node.PublicTransactionList.Head
	for (cur != nil) {
		fmt.Println(cur.Id)
		ans = append(ans, cur.Id)
		cur = cur.Next
	}

	return ans
}

func (net *P2P) ShutDown() {
	for _,n := range net.Nodes {
		n.ShutDown()
	}
}
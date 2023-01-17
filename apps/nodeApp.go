package main

import (
	//"fmt"
	"os"
	"strconv"
	"time"
	"tnguyen/blockchainexample/consensus"
	"tnguyen/blockchainexample/network"
	"tnguyen/blockchainexample/transaction"
)

func main() {
	args := os.Args[1:]

	address := args[0]
	transactionLen := args[1]
	key := args[2]
	m := args[3]
	k := args[4]
	alpha := args[5]
	beta := args[6]

	// fmt.Println("Address is: ", address)
	// fmt.Println("Transactionlen is: ", transactionLen)
	// fmt.Println("Key is: ", key)
	// fmt.Println("m is: ", m)
	// fmt.Println("k is: ", k)
	// fmt.Println("alpha is: ", alpha)
	// fmt.Println("beta is: ", beta)

	options := consensus.SnowballConsensusOptions{}
	options.M, _ = strconv.Atoi(m)
	options.K, _ = strconv.Atoi(k)
	options.Alpha, _ = strconv.Atoi(alpha)
	options.Beta, _ = strconv.Atoi(beta)

	len, _ := strconv.Atoi(transactionLen)
	keyNum, _ := strconv.Atoi(key)
	id, _ := strconv.Atoi(address)

	// create transaction list
	privateList := new(transaction.TransactionList)
	if id <= options.M/2 {
		for len > 0 {
			privateList.Insert(keyNum)
			len--
		}
	} else if id <= options.M*3/4 {
		for len > 0 {
			privateList.Insert(keyNum + 1)
			len--
		}
	} else {
		for len > 0 {
			privateList.Insert(id)
			len--
		}
	}

	publicList := new(transaction.TransactionList)

	node := network.NodeConstructor("tcp", id, publicList, privateList)
	validator := consensus.CreateSnowballConsensus(node, options)
	node.Validator = validator

	node.Start()

	for {
		time.Sleep(1 * time.Second)
	}
}

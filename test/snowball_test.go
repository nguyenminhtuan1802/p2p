package test

import (
	"testing"
	"fmt"
	"tnguyen/blockchainexample/network"
	"tnguyen/blockchainexample/consensus"
)

func TestSnowBall(t *testing.T) {
	nodeNum := 20
	transactionLen := 10
	key := 32

	net := network.CreateP2PNetwork(nodeNum,transactionLen, key)

	options := consensus.SnowballConsensusOptions{}
	options.M = nodeNum
	options.K = 10
	options.Alpha = 3
	options.Beta = 5
	options.TransactionConsensusThreshold = 18

	validator := consensus.CreateSnowballConsensus(net, options)
	answer := make([]int,0)
	for {
		if validator.Update() {
			fmt.Println("DONE")
			answer = net.GetPublicTransactionList()
			break
		}
		validator.Iteration++
	}

	net.ShutDown()	

	if transactionLen != len(answer) {
		t.Fatalf("Wrong public transaction length")
	}

	for _, v := range answer {
		if v != key {
			t.Fatalf("Wrong public transaction value")
		}
	}
}
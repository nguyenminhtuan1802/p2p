package test

import (
	"testing"
	//"fmt"
	//"tnguyen/blockchainexample/network"
	"tnguyen/blockchainexample/consensus"
)

func TestSnowBall(t *testing.T) {
	nodeNum := 200
	transactionLen := 10
	key := 3200 // public transaction data value should converge to this key value
	
	options := consensus.SnowballConsensusOptions{}
	options.M = nodeNum
	options.K = 20
	options.Alpha = 12
	options.Beta = 5
	
	net := CreateP2PNetwork(nodeNum,transactionLen, key, &options)
	net.Run()

	// Check if all nodes have finished running
	for {
		nodeDoneCount := 0
		for _, n := range net.Nodes {
			if n.FinishRunning {
				nodeDoneCount++
			}
		}
		if nodeDoneCount == nodeNum {
			break
		}
	}

	net.ShutDown()
}
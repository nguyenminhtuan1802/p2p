package test

import (
	"testing"
	"fmt"
	"strconv"
	"time"
	//"tnguyen/blockchainexample/network"
	"tnguyen/blockchainexample/consensus"
)

func Test_200Nodes_1Process(t *testing.T) {
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

func Test_200Nodes_200Processes(t *testing.T) {
	nodeNum := 200
	transactionLen := 10
	key := 3200 // public transaction data value should converge to this key value
	
	options := consensus.SnowballConsensusOptions{}
	options.M = nodeNum
	options.K = 20
	options.Alpha = 12
	options.Beta = 5
	
	// Create and launch nodes in network
	net := CreateP2PNetworkDetach(nodeNum,transactionLen, key, &options)

	for {
		count := 0
		for i := 0; i < nodeNum; i++ {
			go net.Client[i].Send("HASFINISHED")
			finished := net.Client[i].ReceiveOnce().(string)
			if (finished == "TRUE") {
				count++
			}
		}
		if (count == nodeNum) {
			break
		}
	}
	fmt.Println("DONE ")

	for i := 0; i < nodeNum; i++ {
		for j := 0; j <transactionLen; j++{
			go net.Client[i].Send("QUERY," + strconv.Itoa(j))
			data := net.Client[i].ReceiveOnce()
			fmt.Println("[NODE ", i+1, "] Transaction[", j, "]: ", data)
		}
	}	
	time.Sleep(2 * time.Second)
	for i := 0; i < nodeNum; i++ {
		net.Client[i].Send("SHUTDOWN")
	}
}
package test

import (
	//"fmt"
	"os/exec"
	"strconv"
	//"sync"
	//"tnguyen/blockchainexample/transaction"
	"tnguyen/blockchainexample/consensus"
	"tnguyen/blockchainexample/network"
	"tnguyen/blockchainexample/utility"

	"time"
)

type P2PDetach struct {
	Client []*network.Client
}

func CreateP2PNetworkDetach(size int, transactionLen int, key int, options *consensus.SnowballConsensusOptions) *P2PDetach {
	net := new(P2PDetach)
	net.Client = make([]*network.Client, 0)
	
	// Launch nodes
	for i := 1; i <= size; i++ {
		cmd := exec.Command("./../bin/nodeApp", strconv.Itoa(i), strconv.Itoa(transactionLen), strconv.Itoa(key),
		 strconv.Itoa(size), strconv.Itoa(options.K), strconv.Itoa(options.Alpha), strconv.Itoa(options.Beta))

		err := cmd.Start()
		if err != nil {
			panic(err)
		}
	}

	// Wait for all servers to start
	time.Sleep(1 * time.Second)

	// Connect nodes in network
	for i := 1; i <= size; i++ {
		client := network.ClientConstructor("tcp", utility.Addresses[i])
		client.Connect()
		for j := 1; j <= size; j++ {
			if j == i {
				continue
			}
			request := "CONNECT," + strconv.Itoa(j)
			client.Send(request)
			time.Sleep(50 * time.Microsecond) // Need small wait because server does not have message queue, so messages might drop if sent too quickly
		}
		net.Client = append(net.Client, client)
	}

	// Start nodes
	for i := 0; i < size; i++ {
		net.Client[i].Send("RUN")
	}

	return net
}
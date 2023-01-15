package consensus

import (
	"fmt"
	"tnguyen/blockchainexample/network"
	"math/rand"
)


type SnowballConsensusOptions struct {
	M int //participants
	K int //sample size
	Alpha int //quarum size
	Beta int //decision threshold 
	TransactionConsensusThreshold int //threshold to update public transaction
}

type TransactionValidator struct {
	OldPreference int
	ConsecutiveSuccess int
	Decided bool
	Seed int
}

type SnowballConsensusValidator struct {
	Net		 		*network.P2P
	TransValidator	map[int]*TransactionValidator
	Options			SnowballConsensusOptions	
	Iteration		int // DEBUG	
}

func CreateSnowballConsensus(net *network.P2P, opt SnowballConsensusOptions) *SnowballConsensusValidator {
	v := new(SnowballConsensusValidator)
	v.Net = net
	v.Options = opt
	v.TransValidator = make(map[int]*TransactionValidator)
	for _, n := range net.Nodes {
		validator := new(TransactionValidator)
		validator.OldPreference = 0
		validator.ConsecutiveSuccess = 0
		validator.Decided = false
		validator.Seed = n.Id		
		v.TransValidator[n.Id] = validator
	}
	v.Iteration = 1
	return v
}

func (validator *SnowballConsensusValidator) Update() bool {
	fmt.Println("ITERATION: ", validator.Iteration)
	decidedCount := 0 // number of nodes have decided its consensus on this transaction
	consensusValue := 0 // the final transaction data (e.g.transaction id)
	for _, n := range validator.Net.Nodes {

		// Find k peers of node n
		r := rand.New(rand.NewSource(int64((*validator.TransValidator[n.Id]).Seed)))

		peersToQuery := make(map[int]*network.Client,0)
		tempK := validator.Options.K

		for tempK > 0 {
			peerIdx := r.Intn(validator.Options.M) + 1
			peer, ok := n.Clients[peerIdx]

			if (ok) {
				_, ok1 := peersToQuery[peerIdx]
				if (!ok1) {
					peersToQuery[peerIdx] = peer
					tempK--
				}
			}
		}

		responseFreqMap := make(map[int]int)

		for _, peer := range peersToQuery {
			// Ask peer
			go peer.Send("QUERY")
			reply := peer.ReceiveOnce()
			val, ok := responseFreqMap[reply]
			if (ok) {
				responseFreqMap[reply] = val + 1
			} else {
				responseFreqMap[reply] = 1
			}
		}

		// Gather reply
		cnt := 0
		perference := 0
		for res, freq := range responseFreqMap {
			if (freq > cnt) {
				perference = res
				cnt = freq
			}
		}

		// Validate based on the current quarum consensus
		decided, val := (*validator.TransValidator[n.Id]).Validate(&n.PrivateTransactionList.GetCurrentTransaction().Id, perference, cnt, &validator.Options)
		consensusValue = val
		if decided  {
			decidedCount++
		}
	}

	fmt.Println("Number of decided nodes is: ", decidedCount)

	NextTrans := false
	Stop := false
	if (decidedCount > validator.Options.TransactionConsensusThreshold) {
		for _, n := range validator.Net.Nodes {
			n.PublicTransactionList.Insert(consensusValue)
			(*validator.TransValidator[n.Id]).Reset()
			if (!n.PrivateTransactionList.MoveToNextTransaction()) {
				Stop = true
			}
			NextTrans = true
		}
	}
	
	if (Stop) {
		fmt.Println("End of transaction list")
	} else if (NextTrans)  {			
		fmt.Println("Move to next transaction")
	}

	return Stop
}

func (v *TransactionValidator) Validate(data *int, perference int, count int, opt *SnowballConsensusOptions) (bool, int) {
	if (count >= opt.Alpha) {
		if (perference == v.OldPreference) {
			v.ConsecutiveSuccess++
		} else {
			v.ConsecutiveSuccess = 1
			v.OldPreference = perference
		}
	} else {
		v.ConsecutiveSuccess = 0
	}

	if (v.ConsecutiveSuccess > opt.Beta) {
		*data = perference
		v.Decided = true
	}
	return v.Decided, perference
}

func (v *TransactionValidator) Reset() {
	v.OldPreference = 0
	v.ConsecutiveSuccess = 0
	v.Decided = false
	v.Seed = v.Seed + 123
}
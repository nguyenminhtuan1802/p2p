package consensus

import (
	"fmt"
	"math/rand"
	"strconv"
	"tnguyen/blockchainexample/network"
	"tnguyen/blockchainexample/transaction"

)

type SnowballConsensusOptions struct {
	M                             int //participants
	K                             int //sample size
	Alpha                         int //quarum size
	Beta                          int //decision threshold
}

type TransactionValidator struct {
	OldPreference      int
	ConsecutiveSuccess int
	Decided            bool
	Seed               int
}

type SnowballConsensusValidator struct {
	Node           *network.Node
	TransValidator *TransactionValidator
	Options        SnowballConsensusOptions
	Iteration      int // DEBUG
}

func CreateSnowballConsensus(node *network.Node, opt SnowballConsensusOptions) *SnowballConsensusValidator {
	v := new(SnowballConsensusValidator)
	v.Node = node
	v.Options = opt

	v.TransValidator = new(TransactionValidator)
	v.TransValidator.OldPreference = 0
	v.TransValidator.ConsecutiveSuccess = 0
	v.TransValidator.Decided = false
	v.TransValidator.Seed = node.Id

	v.Iteration = 1 // TODO: remove
	return v
}

func (validator *SnowballConsensusValidator) Update() bool {
	fmt.Println("[NODE ", validator.Node.Id, "]  ITERATION: ", validator.Iteration)
	consensusValue := 0 // the final transaction data (e.g.transaction id)

	// Find k peers of node n
	r := rand.New(rand.NewSource(int64((*validator.TransValidator).Seed)))

	peersToQuery := make(map[int]*network.Client, 0)
	tempK := validator.Options.K

	for tempK > 0 {
		peerIdx := r.Intn(validator.Options.M) + 1
		peer, ok := validator.Node.Clients[peerIdx]

		if ok {
			_, ok1 := peersToQuery[peerIdx]
			if !ok1 {
				peersToQuery[peerIdx] = peer
				tempK--
			}
		}
	}

	responseFreqMap := make(map[int]int)
	transactionIndex, transaction := validator.Node.GetNextUndecidedTransaction()

	for _, peer := range peersToQuery {
		// Ask peer
		go peer.Send("QUERY," + strconv.Itoa(transactionIndex))
		reply := peer.ReceiveOnce()
		val, ok := responseFreqMap[reply]
		if ok {
			responseFreqMap[reply] = val + 1
		} else {
			responseFreqMap[reply] = 1
		}
	}

	// Gather reply
	cnt := 0
	perference := 0
	for res, freq := range responseFreqMap {
		if freq > cnt {
			perference = res
			cnt = freq
		}
	}

	// Validate based on the current quarum consensus
	decided, val := (*validator.TransValidator).Validate(transaction, perference, cnt, &validator.Options)
	consensusValue = val

	if decided {
		validator.Node.PublicTransactionList.Insert(consensusValue)
		if !validator.Node.PrivateTransactionList.MoveToNextTransaction() {
			fmt.Println("[NODE ", validator.Node.Id, "]  End of transaction list")
			return false
		}
		fmt.Println("[NODE ", validator.Node.Id, "]  Move to next transaction")
		(*validator.TransValidator).Reset()
	}

	validator.Iteration++
	return true
}

func (v *TransactionValidator) Validate(transaction *transaction.Transaction, perference int, count int, opt *SnowballConsensusOptions) (bool, int) {
	if count >= opt.Alpha {
		if perference == v.OldPreference {
			v.ConsecutiveSuccess++
		} else {
			v.ConsecutiveSuccess = 1
			v.OldPreference = perference
		}
	} else {
		v.ConsecutiveSuccess = 0
	}

	if v.ConsecutiveSuccess > opt.Beta {
		transaction.Data = perference
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

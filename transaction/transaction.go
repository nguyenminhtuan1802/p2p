package transaction

type TransactionList struct {
	Transactions	[]*Transaction
	CurrentIdx		int
}

type Transaction struct {
	Data int
}

func (list *TransactionList) Insert(data int) {
	newTransaction := &Transaction{Data: data}

	if (list.Transactions == nil) {
		list.Transactions = make([]*Transaction, 0)
		list.Transactions = append(list.Transactions, newTransaction)
	}

	list.Transactions = append(list.Transactions, newTransaction)
}

func (list *TransactionList) GetTransactionData(index int) *Transaction {
	if (index >= len(list.Transactions)) {
		return nil
	}
	return list.Transactions[index]
}

func (list *TransactionList) MoveToNextTransaction() bool {
	list.CurrentIdx++
	if (list.CurrentIdx >= len(list.Transactions)) {
		return false
	} else {
		return true
	}
}
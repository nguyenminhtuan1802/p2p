package transaction

type TransactionList struct {
	Head *Transaction
	Tail *Transaction
	Current *Transaction
}

type Transaction struct {
	Id int
	Next *Transaction
}

func (list *TransactionList) Insert(id int) {
	newTransaction := &Transaction{Id: id, Next: nil}
	if (list.Head == nil) {
		list.Head = newTransaction
		list.Tail = newTransaction
		list.Current = newTransaction
	} else {
		list.Tail.Next = newTransaction
		list.Tail = newTransaction
	}
}

func (list *TransactionList) GetCurrentTransaction() *Transaction {
	return list.Current
}

func (list *TransactionList) MoveToNextTransaction() bool {
	list.Current = list.Current.Next

	if (list.Current == nil) {
		return false
	} else {
		return true
	}
}
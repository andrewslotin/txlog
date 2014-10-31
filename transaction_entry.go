package txlog

type transactionEntry struct {
	logEntry
	Token     string
	Terminal  bool
	Discarded bool
}

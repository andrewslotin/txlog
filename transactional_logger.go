package txlog

import (
	"math/rand"
	"strconv"
	"sync"
	"time"
)

type transactionalLogger struct {
	sync.Mutex

	writer       Writer
	transactions map[string][]logEntry
	tokenRand    *rand.Rand

	out                  chan transactionEntry
	terminationRequested bool

	workers sync.WaitGroup
}

func NewTransactionalLogger(w Writer) *transactionalLogger {
	return &transactionalLogger{
		writer:    w,
		tokenRand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (l *transactionalLogger) NewEntry() Entry {
	l.Lock()
	defer l.Unlock()

	if l.transactions == nil {
		l.start()
	}

	return l.newTransaction()
}

func (l *transactionalLogger) Terminate() {
	if l.terminationRequested {
		return
	}

	l.terminationRequested = true
	l.workers.Wait()

	l.transactions = nil
	close(l.out)
}

func (l *transactionalLogger) start() {
	// log.Println("Transactional logger is starting")
	l.out = make(chan transactionEntry)
	l.transactions = make(map[string][]logEntry)
	l.terminationRequested = false

	go func() {
		l.workers.Add(1)
		defer l.workers.Done()

		// log.Println("Worker started")

		for {
			select {
			case entry, _ := <-l.out:
				l.Lock()
				// log.Printf("Got an entry [tag:%v, message:%v, terminal:%v]", entry.Tag, entry.Message, entry.Terminal)
				l.transactions[entry.Token] = append(l.transactions[entry.Token], entry.logEntry)

				if entry.Discarded {
					l.discardTransaction(entry.Token)
				}
				if entry.Terminal {
					l.commitTransaction(entry.Token)
				}
				l.Unlock()
			default:
				if l.terminationRequested {
					l.Lock()
					l.commitAllTransactions()
					l.Unlock()
					// log.Println("Worker terminated")
					return
				}
			}
		}
	}()
}

func (l *transactionalLogger) getToken() string {
	return strconv.FormatInt(l.tokenRand.Int63(), 36)
}

func (l *transactionalLogger) newTransaction() *transaction {
	token := l.getToken()
	l.transactions[token] = make([]logEntry, 0)

	return &transaction{
		token: token,
		out:   l.out,
	}
}

func (l *transactionalLogger) commitTransaction(token string) {
	// log.Printf("Committing transaction %s", token)
	entries, ok := l.transactions[token]
	if !ok {
		return
	}
	delete(l.transactions, token)

	entry := aggregateTransactionEntries(entries)
	l.writer.Write(entry)
}

func (l *transactionalLogger) discardTransaction(token string) {
	// log.Printf("Discarding transaction %s", token)
	delete(l.transactions, token)
}

func (l *transactionalLogger) commitAllTransactions() {
	for token, _ := range l.transactions {
		l.commitTransaction(token)
	}
}

func aggregateTransactionEntries(entries []logEntry) logEntry {
	resultingEntry := logEntry{
		Context: NewContext(),
	}

	for _, entry := range entries {
		resultingEntry.Context.Update(entry.Context)
	}

	if len(entries) > 0 {
		commitEntry := entries[len(entries)-1]
		resultingEntry.Tag = commitEntry.Tag
		resultingEntry.Message = commitEntry.Message
	}

	return resultingEntry
}

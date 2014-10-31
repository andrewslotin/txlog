package txlog

import (
	"fmt"
	"sync"
	"testing"

	. "github.com/adjust/gocheck"
)

func TestTransactinalLoggerSuite(t *testing.T) {
	TestingSuiteT(&TransactinalLoggerSuite{}, t)
}

type TransactinalLoggerSuite struct {
	writer   *WriterMock
	txLogger *transactionalLogger
}

func (suite *TransactinalLoggerSuite) SetUpSuite(c *C) {
	suite.writer = NewWriterMock()
	suite.txLogger = NewTransactionalLogger(suite.writer)
}

func (suite *TransactinalLoggerSuite) TearDownTest(c *C) {
	suite.txLogger.Terminate()
	suite.writer.Reset()
}

func (suite *TransactinalLoggerSuite) TestBeginTransaction(c *C) {
	suite.txLogger.NewEntry()
	c.Check(suite.txLogger.transactions, HasLen, 1)

	suite.txLogger.NewEntry()
	c.Check(suite.txLogger.transactions, HasLen, 2)
}

func (suite *TransactinalLoggerSuite) TestCommitTransaction(c *C) {
	tx := suite.txLogger.NewEntry()

	tx.Log("test", Context{"field": "value"})
	tx.Log("test2", Context{"error": false})

	c.Check(suite.writer.Entries, HasLen, 0)

	tx.Commit("test", "commit")

	c.Check(suite.txLogger.transactions, HasLen, 0)

	c.Assert(suite.writer.Entries, HasLen, 1)
	entry := suite.writer.Entries[0]

	c.Check(entry.Tag, Equals, "test")
	c.Check(entry.Message, Equals, "commit")

	c.Assert(entry.Context, HasLen, 2)
	c.Check(entry.Context["field"], Equals, "value")
	c.Check(entry.Context["error"], Equals, false)
}

func (suite *TransactinalLoggerSuite) TestDiscardTransaction(c *C) {
	tx := suite.txLogger.NewEntry()

	tx.Log("test", Context{"field": "value"})
	tx.Log("test2", Context{"error": false})

	c.Check(suite.writer.Entries, HasLen, 0)

	tx.Discard()
	c.Check(suite.txLogger.transactions, HasLen, 0)
	c.Check(suite.writer.Entries, HasLen, 0)

	tx.Commit("test", "after discard")
	c.Check(suite.writer.Entries, HasLen, 0)
}

func (suite *TransactinalLoggerSuite) TestMultipleTransactions(c *C) {
	threadsNum := 10

	wg := sync.WaitGroup{}
	wg.Add(threadsNum)

	for i := 0; i < threadsNum; i++ {
		go func(n int) {
			defer wg.Done()

			tag := fmt.Sprintf("test%d", n)

			tx := suite.txLogger.NewEntry()
			tx.Log(tag, Context{"n": n})
			tx.Commit(tag, fmt.Sprintf("commit message #%d", n))
		}(i)
	}

	wg.Wait()

	c.Check(suite.writer.Entries, HasLen, threadsNum)
}

func (suite *TransactinalLoggerSuite) TestTermination(c *C) {
	openTransactionsNum := 10

	for i := 0; i < openTransactionsNum; i++ {
		tag := fmt.Sprintf("test%d", i)

		suite.txLogger.NewEntry().Log(tag, Context{"n": i})
	}

	suite.txLogger.Terminate()

	c.Check(suite.writer.Entries, HasLen, openTransactionsNum)

	c.Check(suite.txLogger.transactions, HasLen, 0)
	c.Check(func() { suite.txLogger.out <- transactionEntry{} }, PanicMatches, "runtime error: send on closed channel")
}

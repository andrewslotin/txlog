package txlog

import "sync"

// A logger mock to use for testing purposes that satisfies txlog.Writer interface
type WriterMock struct {
	Entries []logEntry
	mutex   sync.Mutex
}

func NewWriterMock() *WriterMock {
	return &WriterMock{
		Entries: make([]logEntry, 0),
	}
}

func (w *WriterMock) Write(entry logEntry) {
	w.mutex.Lock()
	w.Entries = append(w.Entries, entry)
	w.mutex.Unlock()
}

func (w *WriterMock) Reset() {
	w.Entries = make([]logEntry, 0)
}

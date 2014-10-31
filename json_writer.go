package txlog

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"
)

type jsonWriter struct {
	hostname string
	service  string
	encoder  *json.Encoder
	mutex    sync.Mutex
}

func NewJsonLogger(w io.Writer, service string) *jsonWriter {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	return &jsonWriter{
		hostname: hostname,
		service:  service,
		encoder:  json.NewEncoder(w),
	}
}

func (w *jsonWriter) Log(tag, message string, context Context) {
	entry := NewLogEntry(w.hostname, w.service, message)
	entry.Tag = tag
	entry.Context = context

	w.Write(entry)
}

func (w *jsonWriter) Write(entry logEntry) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	entry.Hostname = w.hostname
	entry.Service = w.service
	entry.Timestamp = time.Now()

	w.encoder.Encode(entry)
}

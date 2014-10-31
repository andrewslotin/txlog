package txlog

type transaction struct {
	token  string
	out    chan transactionEntry
	closed bool
}

func (t *transaction) Log(tag string, context Context) {
	t.log(tag, "", context, false, false)
}

func (t *transaction) Commit(tag, message string) {
	t.log(tag, message, nil, true, false)
	t.closed = true
}

func (t *transaction) Discard() {
	t.log("", "", nil, true, true)
	t.closed = true
}

func (t *transaction) log(tag, message string, context Context, terminal bool, discarded bool) {
	if t.closed {
		return
	}

	record := transactionEntry{
		logEntry:  NewLogEntry("", "", message),
		Token:     t.token,
		Terminal:  terminal,
		Discarded: discarded,
	}
	record.Tag = tag
	record.Context = context

	t.write(record)
}

func (t *transaction) write(entry transactionEntry) {
	t.out <- entry
}

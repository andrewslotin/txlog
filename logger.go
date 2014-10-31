package txlog

type Writer interface {
	Write(entry logEntry)
}

type Entry interface {
	Log(tag string, context Context)
	Commit(tag, message string)
	Discard()
}

type Logger interface {
	NewEntry() Entry
	Terminate()
}

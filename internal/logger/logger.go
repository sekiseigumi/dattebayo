package logger

import (
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Source    string
	Message   string
}

type Logger struct {
	entries []LogEntry
	mu      sync.Mutex
	notify  chan struct{}
}

func NewLogger() *Logger {
	return &Logger{
		entries: make([]LogEntry, 0),
		notify:  make(chan struct{}, 100),
	}
}

func (l *Logger) Log(source, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	entry := LogEntry{
		Timestamp: time.Now(),
		Source:    source,
		Message:   message,
	}

	l.entries = append(l.entries, entry)

	select {
	case l.notify <- struct{}{}:
	default:
	}
}

func (l *Logger) Entries() []LogEntry {
	l.mu.Lock()
	defer l.mu.Unlock()

	return append([]LogEntry(nil), l.entries...)
}

func (l *Logger) Notify() <-chan struct{} {
	return l.notify
}

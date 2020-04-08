package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	logging "github.com/op/go-logging"
)

// StructuredLogger struct
type StructuredLogger struct {
	Logger   *logging.Logger
	R        *http.Request
	Excludes []string
}

// LogFields type
type LogFields map[string]interface{}

func contains(list []string, search string) bool {
	set := make(map[string]bool)
	for _, v := range list {
		set[v] = true
	}
	return set[search]
}

// NewLogEntry method
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: l.Logger}
	logFields := LogFields{}

	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Fields = logFields

	if !contains(l.Excludes, r.RequestURI) {
		l.Logger.Infof("[%s] [START] %s %s - %s, %s", logFields["req_id"], logFields["http_method"], logFields["uri"], logFields["remote_addr"], logFields["user_agent"])
	} else {
		logFields["skip"] = true
	}

	return entry
}

// StructuredLoggerEntry struct
type StructuredLoggerEntry struct {
	Logger *logging.Logger
	Fields LogFields
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, intr interface{}) {
	fields := l.Fields

	if fields["skip"] == true {
		return
	}

	fields["resp_status"] = status
	fields["resp_bytes_length"] = bytes
	respElapsed := float64(elapsed.Nanoseconds()) / 1000000.0
	fields["resp_elapsed_ms"] = respElapsed

	msg := fmt.Sprintf("[%s] [END] %s - %v, %dB, %vms", fields["req_id"], fields["uri"], status, bytes, respElapsed)

	if status >= 400 {
		l.Logger.Error(msg)
	} else {
		l.Logger.Info(msg)
	}
}

// Panic method
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	fields := l.Fields

	fields["stack"] = string(stack)
	fields["panic"] = fmt.Sprintf("%+v", v)

	l.Logger.Fatalf("%s, %s", fields["stack"], fields["panic"])
}

// Errorf wraps logger.Errorf
func (l *StructuredLogger) Errorf(format string, args ...interface{}) {
	if reqID := middleware.GetReqID(l.R.Context()); reqID != "" {
		format = fmt.Sprintf("[%s] ", reqID) + format
	}
	l.Logger.Errorf(format, args...)
}

// Error wraps logger.Error
func (l *StructuredLogger) Error(format string) {
	if reqID := middleware.GetReqID(l.R.Context()); reqID != "" {
		format = fmt.Sprintf("[%s] ", reqID) + format
	}
	l.Logger.Error(format)
}

// Debugf wraps logger.Debugf
func (l *StructuredLogger) Debugf(format string, args ...interface{}) {
	if reqID := middleware.GetReqID(l.R.Context()); reqID != "" {
		format = fmt.Sprintf("[%s] ", reqID) + format
	}
	l.Logger.Debugf(format, args...)
}

// Infof wraps logger.Infof
func (l *StructuredLogger) Infof(format string, args ...interface{}) {
	if reqID := middleware.GetReqID(l.R.Context()); reqID != "" {
		format = fmt.Sprintf("[%s] ", reqID) + format
	}
	l.Logger.Infof(format, args...)
}

// Warningf wraps logger.Warningf
func (l *StructuredLogger) Warningf(format string, args ...interface{}) {
	if reqID := middleware.GetReqID(l.R.Context()); reqID != "" {
		format = fmt.Sprintf("[%s] ", reqID) + format
	}
	l.Logger.Warningf(format, args...)
}

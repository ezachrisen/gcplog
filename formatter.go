package gcplog

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/status"
)

const (
	RequestMethod = "requestMethod"
	RequestUrl    = "requestUrl"
	Latency       = "latency"
	// GrpcCode      = "grpcCode"
	// GrpcMessage   = "grpcMessage"
	// GrpcDetails   = "grpcDetails"
	HTTPStatus = "status"
	GrpcStatus = "grpcStatus"
)

type Formatter struct {
	ProjectID string
}

type googleLogEntry struct {
	Message        string          `json:"message"`
	Severity       string          `json:"severity"`
	Additional     logrus.Fields   `json:"additional_info,omitempty"`
	TraceID        string          `json:"logging.googleapis.com/trace,omitempty"`
	Type           string          `json:"@type,omitempty"`
	SourceLocation *sourceLocation `json:"logging.googleapis.com/sourceLocation,omitempty"`
	Request        *Request        `json:"httpRequest,omitempty"`
	GRPCStatus     *GRPCStatus     `json:"grpc,omitempty"`
}

type Request struct {
	RequestMethod string `json:"requestMethod,omitempty"`
	RequestUrl    string `json:"requestUrl,omitempty"`
	Latency       string `json:"latency,omitempty"`
	HTTPStatus    string `json:"status,omitempty"`
}

type GRPCStatus struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Details string `json:"details,omitempty"`
}

type sourceLocation struct {
	File     string `json:"file,omitempty"`
	Line     int    `json:"line,omitempty"`
	Function string `json:"function,omitempty"`
}

const errorType = "type.googleapis.com/google.devtools.clouderrorreporting.v1beta1.ReportedErrorEvent"

var levels = map[logrus.Level]string{
	logrus.InfoLevel:  "INFO",
	logrus.DebugLevel: "DEBUG",
	logrus.TraceLevel: "DEBUG",
	logrus.WarnLevel:  "WARNING",
	logrus.ErrorLevel: "ERROR",
	logrus.PanicLevel: "CRITICAL",
	logrus.FatalLevel: "CRITICAL",
}

func getGCPLevel(level logrus.Level) string {

	levelstring, ok := levels[level]
	if !ok {
		levelstring = "INFO"
	}
	return levelstring
}

// Format logrus output per Google Cloud guidelines.
// See https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry for details.
//
// See the examples for usage.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := getGCPLevel(entry.Level)

	e := googleLogEntry{
		Message:  entry.Message,
		Severity: level, // entry.Level.String(),
		//	Additional: entry.Data,
	}

	if entry.Caller != nil {
		e.SourceLocation = &sourceLocation{
			File:     path.Base(entry.Caller.File),
			Line:     entry.Caller.Line,
			Function: entry.Caller.Function,
		}
	}

	if entry.Level == logrus.ErrorLevel {
		e.Type = errorType
		e.Message = entry.Message + "\n" + string(debug.Stack())
	}

	if entry.Context != nil {
		span := trace.FromContext(entry.Context)
		if span != nil {
			e.TraceID = fmt.Sprintf("projects/%s/traces/%v", f.ProjectID, span.SpanContext().TraceID)
		}
	}

	// Extract a gRPC Status type
	// https://godoc.org/google.golang.org/grpc/status
	if v, ok := entry.Data[GrpcStatus]; ok {
		if err, ok := v.(error); ok {
			if s, ok := status.FromError(err); ok {
				e.GRPCStatus = &GRPCStatus{
					Code:    fmt.Sprintf("%v", s.Code()),
					Message: fmt.Sprintf("%s", s.Message()),
				}
				if len(s.Details()) > 0 {
					e.GRPCStatus.Details = fmt.Sprintf("%v", s.Details())
				}
				delete(entry.Data, GrpcStatus)
			}
		}
	}

	if requestMethod, ok := entry.Data[RequestMethod]; ok && requestMethod != "" {
		e.Request = &Request{
			RequestMethod: fmt.Sprintf("%v", requestMethod),
			RequestUrl:    fmt.Sprintf("%v", entry.Data[RequestUrl]),
			Latency:       fmt.Sprintf("%v", entry.Data[Latency]),
		}
		delete(entry.Data, RequestMethod)
		delete(entry.Data, RequestUrl)
		delete(entry.Data, Latency)
	}

	e.Additional = entry.Data
	serialized, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	return append(serialized, '\n'), nil
}

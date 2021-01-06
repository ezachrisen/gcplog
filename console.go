package applog

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type ConsoleFormatter struct {
	ProjectID string
}

func (f *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	s := fmt.Sprintf("%10s: %s\n", entry.Level.String(), entry.Message)
	return []byte(s), nil

	// e := googleLogEntry{
	// 	Message:    entry.Message,
	// 	Severity:   entry.Level.String(),
	// 	Additional: entry.Data,
	// }

	// if entry.Caller != nil {
	// 	e.SourceLocation = &sourceLocation{
	// 		File:     path.Base(entry.Caller.File),
	// 		Line:     entry.Caller.Line,
	// 		Function: entry.Caller.Function,
	// 	}
	// }

	// if entry.Level == logrus.ErrorLevel {
	// 	e.Type = errorType
	// 	e.Message = entry.Message + "\n" + string(debug.Stack())
	// }

	// if entry.Context != nil {
	// 	span := trace.FromContext(entry.Context)
	// 	if span != nil {
	// 		e.TraceID = fmt.Sprintf("projects/%s/traces/%v", f.ProjectID, span.SpanContext().TraceID)
	// 	}
	// }

	// if entry.Data != nil {
	// 	//		fmt.Printf("entry.Data (%s) = %v\n", e.TraceID, entry.Data)
	// 	if requestMethod, ok := entry.Data[requestMethod]; ok && requestMethod != "" {
	// 		e.HttpRequest = HttpRequest{
	// 			RequestMethod: fmt.Sprintf("%v", requestMethod),
	// 			RequestUrl:    fmt.Sprintf("%v", entry.Data[requestUrl]),
	// 			Latency:       fmt.Sprintf("%s", entry.Data[latency]),
	// 		}
	// 	}

	// 	if code, ok := entry.Data[grpcCode]; ok && code != "" {
	// 		e.GRPCStatus = GRPCStatus{
	// 			Code:    fmt.Sprintf("%v", code),
	// 			Message: fmt.Sprintf("%s", entry.Data[grpcMessage]),
	// 			Details: fmt.Sprintf("%v", entry.Data[grpcDetails]),
	// 		}
	// 	}
	// 	e.Additional = entry.Data
	// }

	// serialized, err := json.Marshal(e)
	// if err != nil {
	// 	return nil, err
	// }
	// return append(serialized, '\n'), nil
}

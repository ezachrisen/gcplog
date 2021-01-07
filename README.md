# GCPLog

GCPLog formats [logrus](https://github.com/sirupsen/logrus) output for Google Cloud Platform:
- Errors are sent to Google Error Reporting with a stacktrace
- Code calling location is formatted with file, line and module
- Trace ID provided in the context is logged appropriately

### Google Cloud Platform Logging Specification
The correct formatting of log entries is specified here: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry.


### Basic Usage

See the examples for more.

```go 
import (
	"github.com/ezachrisen/gcplog"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})

	log.Info("Hello")
	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 1,
	}).Info("My info message here")

	// Output:
	// {"message":"Hello","severity":"INFO"}
	// {"message":"My info message here","severity":"INFO","additional_info":{"animal":"walrus","number":1}}
}
```


### Log Level Mapping

The following mapping is used between log levels:

| Logrus Level | Google Cloud Platform Level |
| --- | --- |
| INFO | INFO |
| DEBUG | DEBUG |
| TRACE | DEBUG |
| WARN, WARNING | WARNING |
| ERROR | ERROR |
| PANIC | CRITICAL |
| FATAL | CRITICAL |
|  | ALERT |
|  | EMERGENCY |
|  | NOTICE |

### Trace IDs

If an OpenCensus trace is present in the context, use it to add it to the log entry. In order to take advantage of this, you must pass the context to logrus when logging:

```go
logrus.WithContext(ctx).Infof("My info here %d", 100)
// Output:
// {"message":"My info here 100","severity":"INFO","logging.googleapis.com/trace":"projects/myproject/traces/31323334353637383961626364656667"}
```

Google Cloud Platform requires that the trace ID be logged with the project name, therefore you must initialize GCPFormatter with the project ID. 


### gRPC Status
When returning an error from a gRPC API handler, the recommended error type to return is a Status from the gRPC Status package (https://godoc.org/google.golang.org/grpc/status). The Status type encapsulates the error code, a message additional details. See https://cloud.google.com/apis/design/errors for more information on the fields. 

Google Cloud Logging does not define gRPC status fields, but we can log it as a custom JSON object so that the individual fields are preserved (rather than logging it as a single string). 

You could do this "manually" with logrus's WithFields, but GCPLog gives you convenience functions for Info, Warn and Error. 

```go

// in my gRPC API handler
... 

if err != nil {
	err = status.Errorf(codes.NotFound, "blah with key '%s' not found: %v", 
	"myid", err)
	gcplog.GrpcInfo(ctx, err)
	return nil, err
}

// Output
{
   "message":"blah with key 'myid' not found",
   "severity":"INFO",
   "logging.googleapis.com/trace":"projects/myproject/traces/31323334353637383961626364656667",
   "logging.googleapis.com/sourceLocation":{
      "file":"example_test.go",
      "line":83,
      "function":"github.com/ezachrisen/gcplog_test.ExampleGrpcStatusConvenience"
   },
   "grpc":{
      "code":"NotFound",
      "message":"blah with key 'myid' not found"
   }
}

```

### Console Formatter
For local debugging and testing, a bare-bones formatter is provided. This is indentended to log to the terminal. 
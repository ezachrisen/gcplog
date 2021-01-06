# GCPLog

GCPLog formats [logrus](https://github.com/sirupsen/logrus) output for Google Cloud Platform:
- Errors are sent to Google Error Reporting with a stacktrace
- Code calling location is formatted with file, line and module
- Trace ID provided in the context is logged appropriately

### Google Cloud Platform Logging Specification
The correct formatting of log entries is specified here: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry.


### Log Level Mapping

The following mapping is used between log levels:

| Logrus Level | GCP Level |
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

```
logrus.WithContext(ctx).Info("My message here")
```

GCP requires that the trace ID be logged with the project name, therefore you must initialize GCPFormatter with the project ID. 

### Structured Request Data
Logrus lets you pass structured fields (key value pairs) to log, like this:

```go
logrus.WithFields(logrus.Fields{
	"latency" : time.Since(start),
	})
```

Google Cloud understands special request-related fields, and are put in the `httpRequest` field in the log entry. 

If you use a Google managed service (such as Cloud Run or App Engine) and provide the trace context to the logger, there's no need to log these fields explicitly. GCP will log a master record for each request that with the HTTP Request information. By connecting individual log entries to the master record via the trace ID, the information is already available. 

However, if you need to log this separately, GCPLog defines constants for a few of them. 

| Key | Sample Use |
| --- | --- | 
| gcpLog.RequestMethod | "GET" |
| gcpLog.RequestURL | "https://blah.com/MyAPIName" |
| gcpLog.Latency | 123ms |
| gcpLog.HTTPStatus | 200 |

You MUST set the gcpLog.RequestMethod field for GCPLog to recognize that you are passing these values, otherwise no `httpRequest` entry will be created.



### GRPC Status
When providing a gRPC status for an API using the standard gRPC status (https://godoc.org/google.golang.org/grpc/status), you can log this directly to GCPLog with the key gcpLog.GrpcStatus. GCPLog will recognize the gRPC status and log each field separately. 

See https://cloud.google.com/apis/design/errors for more information on the fields. 

Example:

```go

logrus.WithField(gcplog.GrpcStatus, status.Errorf(codes.NotFound, "blah with key %s not found", "myid")).Info("Blah")

	// Output:
	// {"message":"Blah","severity":"INFO","grpc":{"code":"NotFound","message":"blah with key myid not found"}}
```


### Basic Usage

See the examples for more.

```go 
import (
	"github.com/ezachrisen/gcplog"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})

	logrus.Info("Hello")

	logrus.WithFields(logrus.Fields{
		"animal": "walrus",
		"number": 1,
	}).Info("My info message here")

	// Output:
	// {"message":"Hello","severity":"INFO"}
	// {"message":"My info message here","severity":"INFO","additional_info":{"animal":"walrus","number":1}}
}
```

### Console Formatter
For local debugging and testing, a bare-bones formatter is provided. This is indentended to log to the terminal. 
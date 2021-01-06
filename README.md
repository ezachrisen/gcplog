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

Google Cloud understands special request-related fields, and are put in the `httpRequest` field in the log entry. GCPLog defines constants for a few of them. 

| Key | Sample Use |
| --- | --- | 
| gcpLog.RequestMethod | "GET" |
| gcpLog.RequestURL | "https://blah.com/MyAPIName" |
| gcpLog.Latency | 123ms |
| gcpLog.HTTPStatus | 200 |

You MUST set the gcpLog.RequestMethod field for GCPLog to recognize that you are passing these values, otherwise no `httpRequest` entry will be created.



### GRPC Fields 
In addition to the request fields understood by GCP, GCPLog also defines a special set of fields related to gRPC requests. These entries are logged in a separate JSON object in the log entry. 

See https://cloud.google.com/apis/design/errors for more information on the fields. 

| Key | Sample Use |
| --- | --- | 
| gcpLog.GrpcCode | 2 |
| gcpLog.GrpcMessage | My message |
| gcpLog.GrpcDetails | []string{"blah", "foo"} |







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


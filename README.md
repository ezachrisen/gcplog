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

Google Cloud understands special request-related fields. GCPLog defines constants for a few of them. 

| Key | Sample Use |
| --- | --- | 
| gcpLog.RequestMethod | "GET" |
| gcpLog.RequestURL | "https://blah.com/MyAPIName" |
| gcpLog.Latency | 123ms |
| gcpLog.HTTPStatus | 200 |





### GRPC Fields 
In addition to the request fields understood by GCP, GCPLog also defines a special set of fields related to gRPC requests. These entries are logged in a separate JSON object in the log entry. 

See https://cloud.google.com/apis/design/errors for more information on the fields. 

| Key |
| --- | 
| gcpLog.GrpcCode |
| gcpLog.GrpcMessage |
| gcpLog.GrpcDetails |







### Basic Usage

```go 
import (
	"github.com/ezachrisen/gcplog"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&applog.Formatter{})
	logrus.Info("Hello")
	// Output: {"message":"Hello","severity":"info"}
}
```

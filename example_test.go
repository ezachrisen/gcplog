package gcplog_test

import (
	"context"
	"os"

	"github.com/ezachrisen/gcplog"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
)

func ExampleBasic() {

	logrus.SetOutput(os.Stdout) // required for testing only
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

func ExampleTraceID() {

	logrus.SetOutput(os.Stdout) // required for testing only
	logrus.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})
	var dummyTraceID [16]byte
	copy(dummyTraceID[:], "123456789abcdefg")
	ctx, _ := trace.StartSpanWithRemoteParent(context.Background(), "main",
		trace.SpanContext{
			TraceID: trace.TraceID(dummyTraceID),
		},
	)

	logrus.WithContext(ctx).Infof("My info here %d", 100)
	logrus.Info("No trace here")
	// Output:
	// {"message":"My info here 100","severity":"INFO","logging.googleapis.com/trace":"projects/myproject/traces/31323334353637383961626364656667"}
	// {"message":"No trace here","severity":"INFO"}
}

func ExampleWithRequestDetails() {

	logrus.SetOutput(os.Stdout) // required for testing only
	logrus.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})

	logrus.WithFields(logrus.Fields{
		gcplog.RequestMethod: "GET",
		gcplog.RequestUrl:    "http://blah.com/myAPI",
		gcplog.Latency:       100,
		"mydata":             "customdata here",
	}).Info("My info message here")

	// Output:
	// {"message":"My info message here","severity":"INFO","additional_info":{"mydata":"customdata here"},"httpRequest":{"requestMethod":"GET","requestUrl":"http://blah.com/myAPI","latency":"100"}}

}

func ExampleWithGRPCDetails() {

	logrus.SetOutput(os.Stdout) // required for testing only
	logrus.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})

	logrus.WithFields(logrus.Fields{
		gcplog.GrpcCode:    2,
		gcplog.GrpcMessage: "My GRPC Message here",
		gcplog.GrpcDetails: []string{"Blah", "foo", "bar"},
	}).Info("My info message here")

	// Output:
	// {"message":"My info message here","severity":"INFO","grpc":{"code":"2","message":"My GRPC Message here","details":"[Blah foo bar]"}}
}

func ExampleError() {

	logrus.SetOutput(os.Stdout) // required for testing only
	logrus.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})

	logrus.Error("NOOOOOO!")
	// No output shown; run this separately to see the stacktrace
}

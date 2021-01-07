package gcplog_test

import (
	"context"
	"os"

	"github.com/ezachrisen/gcplog"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func ExampleError() {

	logrus.SetOutput(os.Stdout) // required for testing only
	logrus.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})

	logrus.Error("NOOOOOO!")
	// No output shown; run this separately to see the stacktrace
}

func ExampleGrpcStatus() {

	logrus.SetOutput(os.Stdout) // required for testing only
	logrus.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})

	logrus.WithField(gcplog.GrpcStatus, status.Errorf(codes.NotFound, "blah with key %s not found", "myid")).Info("Blah")
	// Output:
	// {"message":"Blah","severity":"INFO","grpc":{"code":"NotFound","message":"blah with key myid not found"}}
}

func ExampleGrpcStatusConvenience() {

	logrus.SetOutput(os.Stdout) // required for testing only
	logrus.SetFormatter(&gcplog.Formatter{ProjectID: "myproject"})

	var dummyTraceID [16]byte
	copy(dummyTraceID[:], "123456789abcdefg")
	ctx, _ := trace.StartSpanWithRemoteParent(context.Background(), "main",
		trace.SpanContext{
			TraceID: trace.TraceID(dummyTraceID),
		},
	)

	gcplog.GrpcInfo(ctx, status.Errorf(codes.NotFound, "blah with key '%s' not found", "myid"))
	// Output:
	// {"message":"blah with key 'myid' not found","severity":"INFO","logging.googleapis.com/trace":"projects/myproject/traces/31323334353637383961626364656667","grpc":{"code":"NotFound","message":"blah with key 'myid' not found"}}
}

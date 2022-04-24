package lambda

import (
	"context"
	"fmt"
	"path"
	"runtime"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/sirupsen/logrus"
)

// Klogger is a container for logrus logger
type Klogger struct {
	*logrus.Entry
}

func packageLogger(baseLogger *logrus.Logger, level logrus.Level) *Klogger {
	baseLogger.SetLevel(level)
	baseLogger.SetReportCaller(true)
	baseLogger.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	})
	return &Klogger{baseLogger.WithFields(logrus.Fields{})}
}

// LambdaLoggerCustom is an instance of Klogger with default formatting and custom logger
func LambdaLoggerCustom(level logrus.Level, baseLogger *logrus.Logger) *Klogger {
	return packageLogger(baseLogger, level)
}

// LambdaLogger is an instance of Klogger with default formatting and default logger
func LambdaLogger(level logrus.Level) *Klogger {
	baseLogger := logrus.New()
	return packageLogger(baseLogger, level)
}

// WithLambdaContext sets lambda context for the logger and logs lambda_request_id each time
func (o *Klogger) WithLambdaContext(ctx context.Context) *Klogger {
	var requestID string

	if lc, ok := lambdacontext.FromContext(ctx); ok {
		requestID = lc.AwsRequestID
	} else {
		o.WithField("context", ctx).Warn("Fail to extract lambda context")
		requestID = "N/A"
	}
	return &Klogger{o.
		WithField("lambda_request_id", requestID),
	}
}

// CheckErrorVerbose should be used in all functions as a deferred method call. The parameter
// messageWhenErrorNotNil will be used when the error is NOT nil.  The parameter
// messageWhenErrorIsNil will be used when the error is not nil.
// Example: defer logger.CheckErrorVerbose(err, "error is not nil", "error is nil")
func (o *Klogger) CheckErrorVerbose(
	err *error,
	messageWhenErrorNotNil string,
	messageWhenErrorIsNil string,
) {
	if err != nil {
		o.WithField("error", err).Error(messageWhenErrorNotNil)
	} else {
		o.Infof(messageWhenErrorIsNil)
	}
}

// CheckError should be used in all functions as a deferred method call. The parameter
// messageWhenErrorNotNil will be used when the error is NOT nil.
// Example: defer logger.CheckErrorVerbose(err, "error is not nil")
func (o *Klogger) CheckError(
	err *error,
	errorMessages ...string,
) {
	if err != nil {
		if len(errorMessages) > 0 {
			for _, m := range errorMessages {
				o.WithField("error", err).Error(m)
			}
		} else {
			o.WithField("error", err).Error("error is NOT nil in CheckError")
		}
	}
}

// IsNil should be used in conditions like if statements so that we can trap exactly where
// an error is not nil but should be and trap this along with any debugging information
// Example #1: if IsNil(err) {}
// Example #2: if IsNil(err, "message one", "message two") {}
func (o *Klogger) IsNil(
	err *error,
	errorMessages ...string,
) (rv bool) {
	rv = err == nil
	if !rv {
		if len(errorMessages) > 0 {
			for _, m := range errorMessages {
				o.WithField("error", err).Error(m)
			}
		} else {
			o.WithField("error", err).Error("error is NOT nil in condition")
		}
	}

	return rv
}

// NotNil should be used in conditions like if statements so that we can trap exactly where
// an error is not nil but should be and trap this along with any debugging information
// Example #1: if NotNil(err) {}
// Example #2: if NotNil(err, "message one", "message two") {}
func (o *Klogger) NotNil(
	err *error,
	errorMessages ...string,
) (rv bool) {
	rv = err != nil
	if !rv {
		if len(errorMessages) > 0 {
			for _, m := range errorMessages {
				o.WithField("error", err).Error(m)
			}
		} else {
			o.WithField("error", err).Error("error is nil in condition")
		}
	}

	return rv
}

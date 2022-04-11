package domain_test

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus/hooks/test"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/google/uuid"

	. "github.com/kraneware/kws/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
)

type mockLambdaContext struct{}

func (ctx mockLambdaContext) Deadline() (deadline time.Time, ok bool) {
	return deadline, ok
}

func (ctx mockLambdaContext) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (ctx mockLambdaContext) Err() error {
	return context.DeadlineExceeded
}

func (ctx mockLambdaContext) Value(key interface{}) interface{} {
	return &lambdacontext.LambdaContext{
		AwsRequestID: uuid.New().String(),
	}
}

type mockContext struct{}

func (ctx mockContext) Deadline() (deadline time.Time, ok bool) {
	return deadline, ok
}

func (ctx mockContext) Done() <-chan struct{} {
	ch := make(chan struct{})
	close(ch)
	return ch
}

func (ctx mockContext) Err() error {
	return context.DeadlineExceeded
}

func (ctx mockContext) Value(key interface{}) interface{} {
	return ""
}

var _ = Describe("LambdaLogger Tests", func() {
	baseLogger := LambdaLogger(logrus.InfoLevel)
	Context("Lambda Context", func() {
		It("should log with lambda context", func() {
			lambdaLogger := baseLogger.WithLambdaContext(mockLambdaContext{})
			lambdaLogger.Infof("test logging with context")
		})

		It("should log without lambda context", func() {
			lambdaLogger := baseLogger.WithLambdaContext(mockContext{})
			lambdaLogger.Infof("test logging with no context")
		})

		It("should log error with defer", func() {
			logger, hook := test.NewNullLogger()
			testLogger := LambdaLoggerCustom(logrus.InfoLevel, logger)

			f := func() {
				deferredError := fmt.Errorf("deferred error")
				defer testLogger.CheckErrorVerbose(
					&deferredError,
					"deferred logging for error",
					"no error",
				)
				testLogger.Infof("this logging should come first")
			}

			f()
			Expect(len(hook.Entries)).Should(Equal(2))
			Expect(hook.LastEntry().Message).Should(Equal("deferred logging for error"))
			Expect(hook.AllEntries()[0].Message).Should(Equal("this logging should come first"))
			Expect(logrus.ErrorLevel).Should(Equal(hook.LastEntry().Level))
			hook.Reset()
			Expect(hook.LastEntry()).Should(BeNil())
		})

		It("should NOT log error with defer", func() {
			logger, hook := test.NewNullLogger()
			testLogger := LambdaLoggerCustom(logrus.InfoLevel, logger)

			f := func() {
				defer testLogger.CheckErrorVerbose(
					nil,
					"deferred logging for error",
					"deferred logging for NO error",
				)
				testLogger.Infof("this logging should come first")
			}

			f()
			Expect(len(hook.Entries)).Should(Equal(2))
			Expect(hook.LastEntry().Message).Should(Equal("deferred logging for NO error"))
			Expect(hook.AllEntries()[0].Message).Should(Equal("this logging should come first"))
			Expect(logrus.InfoLevel).Should(Equal(hook.LastEntry().Level))
			hook.Reset()
			Expect(hook.LastEntry()).Should(BeNil())
		})

		It("should log error with defer", func() {
			logger, hook := test.NewNullLogger()
			testLogger := LambdaLoggerCustom(logrus.InfoLevel, logger)

			f := func() {
				deferredError := fmt.Errorf("deferred error")
				defer testLogger.CheckError(
					&deferredError,
					"deferred logging for error",
				)
				testLogger.Infof("this logging should come first")
			}

			f()
			Expect(len(hook.Entries)).Should(Equal(2))
			Expect(hook.LastEntry().Message).Should(Equal("deferred logging for error"))
			Expect(hook.AllEntries()[0].Message).Should(Equal("this logging should come first"))
			Expect(logrus.ErrorLevel).Should(Equal(hook.LastEntry().Level))
			hook.Reset()
			Expect(hook.LastEntry()).Should(BeNil())
		})

		It("should log error with defer and default message", func() {
			logger, hook := test.NewNullLogger()
			testLogger := LambdaLoggerCustom(logrus.InfoLevel, logger)

			f := func() {
				deferredError := fmt.Errorf("deferred error")
				defer testLogger.CheckError(&deferredError)
				testLogger.Infof("this logging should come first")
			}

			f()
			Expect(len(hook.Entries)).Should(Equal(2))
			Expect(hook.LastEntry().Message).Should(Equal("error is NOT nil in CheckError"))
			Expect(hook.AllEntries()[0].Message).Should(Equal("this logging should come first"))
			Expect(logrus.ErrorLevel).Should(Equal(hook.LastEntry().Level))
			hook.Reset()
			Expect(hook.LastEntry()).Should(BeNil())
		})

		It("should log nil error in conditional", func() {
			logger, hook := test.NewNullLogger()
			testLogger := LambdaLoggerCustom(logrus.InfoLevel, logger)

			if testLogger.NotNil(nil) {
				testLogger.Infof("this logging should not show")
			}
			Expect(len(hook.Entries)).Should(Equal(1))
			Expect(hook.AllEntries()[0].Message).Should(Equal("error is nil in condition"))

			hook.Reset()
			if testLogger.NotNil(
				nil,
				"this logging should come first",
				"this logging should come second",
			) {
				testLogger.Infof("this logging should not show up")
			}
			Expect(len(hook.Entries)).Should(Equal(2))
			Expect(hook.AllEntries()[0].Message).Should(Equal("this logging should come first"))
			Expect(hook.AllEntries()[1].Message).Should(Equal("this logging should come second"))

			hook.Reset()
			newError := fmt.Errorf("error")
			if testLogger.NotNil(
				&newError,
				"this first logging should come show",
			) {
				testLogger.Infof("this is the only logging that should show")
			}
			Expect(len(hook.Entries)).Should(Equal(1))
			Expect(hook.AllEntries()[0].Message).Should(Equal("this is the only logging that should show"))
		})

		It("should log error in conditional", func() {
			logger, hook := test.NewNullLogger()
			testLogger := LambdaLoggerCustom(logrus.InfoLevel, logger)

			newError := fmt.Errorf("error")
			if testLogger.IsNil(&newError) {
				testLogger.Infof("this logging should not show")
			}
			Expect(len(hook.Entries)).Should(Equal(1))
			Expect(hook.AllEntries()[0].Message).Should(Equal("error is NOT nil in condition"))

			hook.Reset()
			if testLogger.IsNil(
				&newError,
				"this logging should come first",
				"this logging should come second",
			) {
				testLogger.Infof("this logging should not show up")
			}
			Expect(len(hook.Entries)).Should(Equal(2))
			Expect(hook.AllEntries()[0].Message).Should(Equal("this logging should come first"))
			Expect(hook.AllEntries()[1].Message).Should(Equal("this logging should come second"))

			hook.Reset()
			if testLogger.IsNil(
				nil,
				"this first logging should come show",
			) {
				testLogger.Infof("this is the only logging that should show")
			}
			Expect(len(hook.Entries)).Should(Equal(1))
			Expect(hook.AllEntries()[0].Message).Should(Equal("this is the only logging that should show"))
		})
	})
})

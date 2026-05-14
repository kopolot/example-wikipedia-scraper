package testutils

import (
	"io"
	"os"

	"github.com/stretchr/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, keysAndValues ...interface{})           {}
func (m *MockLogger) Info(msg string, keysAndValues ...interface{})            {}
func (m *MockLogger) Warn(msg string, keysAndValues ...interface{})            {}
func (m *MockLogger) Error(msg string, keysAndValues ...interface{})           {}
func (m *MockLogger) Log(level int8, msg string, keysAndValues ...interface{}) {}
func (m *MockLogger) Fatal(msg string, keysAndValues ...interface{})           {}
func (m *MockLogger) GetLogWriter() io.Writer {
	f, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0644)
	if err != nil {
		return nil
	}
	return f
}

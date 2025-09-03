package phocus_comms

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPort is a mock implementation of the Port interface
type MockPort struct {
	mock.Mock
}

func (m *MockPort) Open() (Port, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(Port), args.Error(1)
}

func (m *MockPort) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockPort) Read(timeout time.Duration) (string, error) {
	args := m.Called(timeout)
	return args.String(0), args.Error(1)
}

func (m *MockPort) Write(data string) (int, error) {
	args := m.Called(data)
	return args.Int(0), args.Error(1)
}

func TestPort_Open(t *testing.T) {
	mockPort := new(MockPort)
	mockPort.On("Open").Return(mockPort, nil)

	port, err := mockPort.Open()
	assert.NoError(t, err)
	assert.Equal(t, mockPort, port)

	mockPort.AssertExpectations(t)
}

func TestPort_Close(t *testing.T) {
	mockPort := new(MockPort)
	mockPort.On("Close").Return(nil)

	err := mockPort.Close()
	assert.NoError(t, err)

	mockPort.AssertExpectations(t)
}

func TestPort_Read(t *testing.T) {
	mockPort := new(MockPort)
	expectedData := "test data"
	mockPort.On("Read", mock.AnythingOfType("time.Duration")).Return(expectedData, nil)

	data, err := mockPort.Read(1 * time.Second)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)

	mockPort.AssertExpectations(t)
}

func TestPort_Write(t *testing.T) {
	mockPort := new(MockPort)
	expectedBytes := 9
	mockPort.On("Write", "test data").Return(expectedBytes, nil)

	bytesWritten, err := mockPort.Write("test data")
	assert.NoError(t, err)
	assert.Equal(t, expectedBytes, bytesWritten)

	mockPort.AssertExpectations(t)
}

func TestPort_Open_Error(t *testing.T) {
	mockPort := new(MockPort)
	mockPort.On("Open").Return(nil, errors.New("open error"))

	port, err := mockPort.Open()
	assert.Error(t, err)
	assert.Nil(t, port)

	mockPort.AssertExpectations(t)
}

func TestPort_Close_Error(t *testing.T) {
	mockPort := new(MockPort)
	mockPort.On("Close").Return(errors.New("close error"))

	err := mockPort.Close()
	assert.Error(t, err)

	mockPort.AssertExpectations(t)
}

func TestPort_Read_Error(t *testing.T) {
	mockPort := new(MockPort)
	mockPort.On("Read", mock.AnythingOfType("time.Duration")).Return("", errors.New("read error"))

	data, err := mockPort.Read(1 * time.Second)
	assert.Error(t, err)
	assert.Empty(t, data)

	mockPort.AssertExpectations(t)
}

func TestPort_Write_Error(t *testing.T) {
	mockPort := new(MockPort)
	mockPort.On("Write", "test data").Return(0, errors.New("write error"))

	bytesWritten, err := mockPort.Write("test data")
	assert.Error(t, err)
	assert.Equal(t, 0, bytesWritten)

	mockPort.AssertExpectations(t)
}

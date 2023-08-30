package phocus_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // for generating UUIDs for commands
	"github.com/gorilla/websocket"
	messages "github.com/wolffshots/phocus/v2/messages"

	"github.com/stretchr/testify/assert"
)

func TestAddQPGSnMessages(t *testing.T) {
	assert.Equal(t, 1, len(Queue))
	assert.Equal(t, "QID", Queue[0].Command)
	err := AddQPGSnMessages(0)
	assert.Equal(t, 3, len(Queue))
	assert.Equal(t, "QID", Queue[0].Command)
	assert.Equal(t, "QPGS1", Queue[1].Command)
	assert.Equal(t, "QPGS2", Queue[2].Command)
	assert.NoError(t, err)
	err = AddQPGSnMessages(0) // shouldn't add to the queue since there is already over 2
	assert.Equal(t, 3, len(Queue))
	assert.Equal(t, errors.New("queue too long"), err)
}

func TestPostMessage(t *testing.T) {
	router := SetupRouter(gin.TestMode)

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	Queue = make([]messages.Message, 0)

	// test first insertion
	want := []messages.Message{
		{ID: qidUUID2, Command: "QID2", Payload: ""},
	}
	body, err := json.Marshal(messages.Message{ID: qidUUID2, Command: "QID2", Payload: ""})
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, want, Queue)

	// test second insertion
	want = []messages.Message{
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	body, err = json.Marshal(messages.Message{ID: qidUUID3, Command: "QID3", Payload: ""})
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, want, Queue)

	// test third insertion
	want = []messages.Message{
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
		{ID: qidUUID1, Command: "QID1", Payload: ""},
	}
	body, err = json.Marshal(messages.Message{ID: qidUUID1, Command: "QID1", Payload: ""})
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, want, Queue)

	// test invalid message insertion
	want = []messages.Message{
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
		{ID: qidUUID1, Command: "QID1", Payload: ""},
	}
	body, err = json.Marshal(nil)
	assert.NoError(t, err)
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, want, Queue)

	// test too many insertions (MAX_QUEUE_LENGTH)
	for i := 4; i <= MAX_QUEUE_LENGTH; i++ {
		body, err = json.Marshal(messages.Message{ID: uuid.New(), Command: fmt.Sprintf("QID%d", i), Payload: ""})
		assert.NoError(t, err)
		w = httptest.NewRecorder()
		req, err = http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)
	}
	assert.Equal(t, MAX_QUEUE_LENGTH, len(Queue))
	body, err = json.Marshal(messages.Message{ID: uuid.New(), Command: fmt.Sprintf("QID%d", 51), Payload: ""})
	assert.NoError(t, err)
	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInsufficientStorage, w.Code)
	assert.Equal(t, MAX_QUEUE_LENGTH, len(Queue)) // should have prevented that insertion
}

func TestGetQueue(t *testing.T) {
	router := SetupRouter(gin.TestMode)

	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()

	Queue = []messages.Message{
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}

	want := []messages.Message{
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/queue", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var actualBody []messages.Message
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, want, actualBody)
	assert.Equal(t, want, Queue)
}

func TestGetHealth(t *testing.T) {
	router := SetupRouter(gin.TestMode)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "UP", w.Body.String())
}

func TestSetAndGetLast(t *testing.T) {
	router := SetupRouter(gin.TestMode)

	SetLast((*messages.QPGSnResponse)(nil))

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/last", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "null", w.Body.String())

	// test with realistic response
	input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf2\x2d\r"
	actual, err := messages.InterpretQPGSn(input, 1)
	assert.Equal(t, err, nil)
	SetLast(actual)

	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/last", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "{\"InverterNumber\":1,\"OtherUnits\":true,\"SerialNumber\":\"92932004102443\",\"OperationMode\":\"Off-grid\",\"FaultCode\":\"\",\"ACInputVoltage\":\"237.0\",\"ACInputFrequency\":\"50.01\",\"ACOutputVoltage\":\"000.0\",\"ACOutputFrequency\":\"00.00\",\"ACOutputApparentPower\":\"0483\",\"ACOutputActivePower\":\"0387\",\"PercentageOfNominalOutputPower\":\"009\",\"BatteryVoltage\":\"51.1\",\"BatteryChargingCurrent\":\"000\",\"BatteryStateOfCharge\":\"069\",\"PVInputVoltage\":\"020.4\",\"TotalChargingCurrent\":\"000\",\"TotalACOutputApparentPower\":\"00942\",\"TotalACOutputActivePower\":\"00792\",\"TotalPercentageOfNominalOutputPower\":\"007\",\"InverterStatus\":{\"MPPT\":\"off\",\"ACCharging\":\"off\",\"SolarCharging\":\"off\",\"BatteryStatus\":\"Battery voltage normal\",\"ACInput\":\"connected\",\"ACOutput\":\"on\",\"Reserved\":\"0\"},\"ACOutputMode\":\"Parallel output\",\"BatteryChargerSourcePriority\":\"Solar first\",\"MaxChargingCurrentSet\":\"060\",\"MaxChargingCurrentPossible\":\"080\",\"MaxACChargingCurrentSet\":\"10\",\"PVInputCurrent\":\"00.0\",\"BatteryDischargeCurrent\":\"006\",\"Checksum\":\"0xf22d\"}", w.Body.String())

	// test with realistic response
	input = "(1 92932004102543 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf2\x2d\r"
	actual, err = messages.InterpretQPGSn(input, 2)
	assert.Equal(t, err, nil)
	SetLast(actual)

	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/last", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, "{\"InverterNumber\":2,\"OtherUnits\":true,\"SerialNumber\":\"92932004102543\",\"OperationMode\":\"Off-grid\",\"FaultCode\":\"\",\"ACInputVoltage\":\"237.0\",\"ACInputFrequency\":\"50.01\",\"ACOutputVoltage\":\"000.0\",\"ACOutputFrequency\":\"00.00\",\"ACOutputApparentPower\":\"0483\",\"ACOutputActivePower\":\"0387\",\"PercentageOfNominalOutputPower\":\"009\",\"BatteryVoltage\":\"51.1\",\"BatteryChargingCurrent\":\"000\",\"BatteryStateOfCharge\":\"069\",\"PVInputVoltage\":\"020.4\",\"TotalChargingCurrent\":\"000\",\"TotalACOutputApparentPower\":\"00942\",\"TotalACOutputActivePower\":\"00792\",\"TotalPercentageOfNominalOutputPower\":\"007\",\"InverterStatus\":{\"MPPT\":\"off\",\"ACCharging\":\"off\",\"SolarCharging\":\"off\",\"BatteryStatus\":\"Battery voltage normal\",\"ACInput\":\"connected\",\"ACOutput\":\"on\",\"Reserved\":\"0\"},\"ACOutputMode\":\"Parallel output\",\"BatteryChargerSourcePriority\":\"Solar first\",\"MaxChargingCurrentSet\":\"060\",\"MaxChargingCurrentPossible\":\"080\",\"MaxACChargingCurrentSet\":\"10\",\"PVInputCurrent\":\"00.0\",\"BatteryDischargeCurrent\":\"006\",\"Checksum\":\"0xf22d\"}", w.Body.String())

}

func TestGetLastStateOfCharge(t *testing.T) {
	router := SetupRouter(gin.TestMode)

	// test with empty LastQPGSResponse (like pre first request)
	SetLast(nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/last/soc", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"BatteryStateOfCharge\":\"null\"}", w.Body.String())

	// test with relistic response
	input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf2\xaa\r"
	actual, err := messages.InterpretQPGSn(input, 3)
	assert.Equal(t, err, nil)
	SetLast(actual)

	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/last/soc", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "{\"BatteryStateOfCharge\":\"069\"}", w.Body.String())
}

func TestGetMessage(t *testing.T) {
	router := SetupRouter(gin.TestMode)

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	qidUUID4 := uuid.New()

	Queue = []messages.Message{
		{ID: qidUUID1, Command: "QID1", Payload: ""},
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	var actualBody messages.Message

	// test next get
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/queue/next", nil)
	router.ServeHTTP(w, req)
	want := messages.Message{ID: qidUUID1, Command: "QID1", Payload: ""}
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, want, actualBody)

	// test first get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/queue/%s", qidUUID2), nil)
	router.ServeHTTP(w, req)
	want = messages.Message{ID: qidUUID2, Command: "QID2", Payload: ""}
	err = json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, want, actualBody)

	// test second get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/queue/%s", qidUUID1), nil)
	router.ServeHTTP(w, req)
	want = messages.Message{ID: qidUUID1, Command: "QID1", Payload: ""}
	err = json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, want, actualBody)

	// test third get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/queue/%s", qidUUID3), nil)
	router.ServeHTTP(w, req)
	want = messages.Message{ID: qidUUID3, Command: "QID3", Payload: ""}
	err = json.Unmarshal(w.Body.Bytes(), &actualBody)
	assert.Equal(t, err, nil)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, want, actualBody)

	// test invalid uuid get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/queue/%s", "IAMNOTAVALIDUUID"), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code) // doesn't convert to uuid on get request

	// test not found uuid get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/queue/%s", qidUUID4), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteQueue(t *testing.T) {
	router := SetupRouter(gin.TestMode)

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	Queue = []messages.Message{
		{ID: qidUUID1, Command: "QID1", Payload: ""},
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	want := make([]messages.Message, 0)

	// test filled case
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/queue", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, want, Queue)

	// test empty case
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req) // ensure no deadlock
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, want, Queue)
}

func TestDeleteMessage(t *testing.T) {
	router := SetupRouter(gin.TestMode)

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	qidUUID4 := uuid.New()
	Queue = []messages.Message{
		{ID: qidUUID1, Command: "QID1", Payload: ""},
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}

	// test non existent deletion
	want := []messages.Message{
		{ID: qidUUID1, Command: "QID1", Payload: ""},
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID4), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, want, Queue)

	// test first deletion
	want = []messages.Message{
		{ID: qidUUID1, Command: "QID1", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID2), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, want, Queue)

	// test second deletion
	want = []messages.Message{
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID1), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, want, Queue)

	// test third deletion
	want = make([]messages.Message, 0)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID3), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, want, Queue)

	// test invalid deletion
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID3), nil)
	router.ServeHTTP(w, req) // ensure no deadlock and test negative case
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, want, Queue) // assert that it hasn't changed
}

func TestQueueQPGSn(t *testing.T) {
	// Start the adder in a goroutine
	go QueueQPGSn(100, 5)

	// Wait for a specific duration to allow the server to start
	time.Sleep(51 * time.Millisecond)

	assert.Equal(t, 1, len(Queue))
}

func TestLastAndLastWS(t *testing.T) {
	// Create a test router
	router := SetupRouter(gin.TestMode)

	// Create a test HTTP server
	ts := httptest.NewServer(router)
	defer ts.Close()

	// Test regular GET request to "/last"
	resp, err := http.Get(ts.URL + "/last")
	if err != nil {
		t.Fatalf("http.Get error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status %d but got %d", http.StatusOK, resp.StatusCode)
	}

	// Test WebSocket functionality
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial("ws"+ts.URL[4:]+"/last-ws", nil)
	if err != nil {
		t.Fatalf("websocket.Dial error: %v", err)
	}
	defer conn.Close()

	input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf3\xab\r"
	actual, err := messages.InterpretQPGSn(input, 4)
	assert.Equal(t, err, nil)
	SetLast(actual)
	// Listen to and verify the WebSocket message
	_, _, err = conn.ReadMessage() // read first message and ignore it
	if err != nil {
		t.Fatalf("conn.ReadMessage error: %v", err)
	}
	var receivedMessage []byte
	for messageCount := 0; messageCount < 10; messageCount++ {
		_, receivedMessage, err = conn.ReadMessage()
	}

	assert.Equal(t, []byte(messages.EncodeQPGSn(actual)), receivedMessage)
	assert.Equal(t, nil, err)
}

func TestUpgraderError(t *testing.T) {
	// Create a test router
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return false
		},
	}
	router := SetupRouter(gin.TestMode)
	ts := httptest.NewServer(router)
	defer ts.Close()
	dialer := websocket.DefaultDialer
	_, _, err := dialer.Dial("ws"+ts.URL[4:]+"/last-ws", nil)
	assert.Error(t, errors.New("bad handshake"), err) // because it couldn't be opened
}

func TestWriteError(t *testing.T) {
	// Create a test router
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		WriteBufferSize: 1,
	}
	router := SetupRouter(gin.TestMode)
	ts := httptest.NewServer(router)
	defer ts.Close()
	dialer := websocket.DefaultDialer
	_, _, err := dialer.Dial("ws"+ts.URL[4:]+"/last-ws", nil)
	assert.Error(t, errors.New("bad handshake"), err) // because it couldn't be written
}

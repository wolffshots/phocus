package phocus_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid" // for generating UUIDs for commands
	messages "github.com/wolffshots/phocus/v2/messages"

	"github.com/stretchr/testify/assert"
)

func TestPostMessage(t *testing.T) {
	router := SetupRouter()

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	Queue = make([]messages.Message, 0)

	// test first insertion
	want := []messages.Message{
		{ID: qidUUID2, Command: "QID2", Payload: ""},
	}
	body, _ := json.Marshal(messages.Message{ID: qidUUID2, Command: "QID2", Payload: ""})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, want, Queue)

	// test second insertion
	want = []messages.Message{
		{ID: qidUUID2, Command: "QID2", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	body, _ = json.Marshal(messages.Message{ID: qidUUID3, Command: "QID3", Payload: ""})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
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
	body, _ = json.Marshal(messages.Message{ID: qidUUID1, Command: "QID1", Payload: ""})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
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
	body, err := json.Marshal(nil)
	assert.Equal(t, nil, err)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/queue", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, want, Queue)
}

func TestGetQueue(t *testing.T) {
	router := SetupRouter()

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
	router := SetupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "UP", w.Body.String())
}

func TestGetLast(t *testing.T) {
	router := SetupRouter()

	// test with empty LastQPGSResponse (like pre first request)
	LastQPGSResponse = (*messages.QPGSnResponse)(nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/last", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "null", w.Body.String())

	// test with relistic response
	input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf2\x2d\r"
	actual, err := messages.NewQPGSnResponse(input, 1)
	assert.Equal(t, err, nil)
	LastQPGSResponse = actual

	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/last", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "{\"InverterNumber\":1,\"OtherUnits\":true,\"SerialNumber\":\"92932004102443\",\"OperationMode\":\"Off-grid\",\"FaultCode\":\"\",\"ACInputVoltage\":\"237.0\",\"ACInputFrequency\":\"50.01\",\"ACOutputVoltage\":\"000.0\",\"ACOutputFrequency\":\"00.00\",\"ACOutputApparentPower\":\"0483\",\"ACOutputActivePower\":\"0387\",\"PercentageOfNominalOutputPower\":\"009\",\"BatteryVoltage\":\"51.1\",\"BatteryChargingCurrent\":\"000\",\"BatteryStateOfCharge\":\"069\",\"PVInputVoltage\":\"020.4\",\"TotalChargingCurrent\":\"000\",\"TotalACOutputApparentPower\":\"00942\",\"TotalACOutputActivePower\":\"00792\",\"TotalPercentageOfNominalOutputPower\":\"007\",\"InverterStatus\":{\"MPPT\":\"off\",\"ACCharging\":\"off\",\"SolarCharging\":\"off\",\"BatteryStatus\":\"Battery voltage normal\",\"ACInput\":\"connected\",\"ACOutput\":\"on\",\"Reserved\":\"0\"},\"ACOutputMode\":\"Parallel output\",\"BatteryChargerSourcePriority\":\"Solar first\",\"MaxChargingCurrentSet\":\"060\",\"MaxChargingCurrentPossible\":\"080\",\"MaxACChargingCurrentSet\":\"10\",\"PVInputCurrent\":\"00.0\",\"BatteryDischargeCurrent\":\"006\",\"Checksum\":\"0xf22d\"}", w.Body.String())
}

func TestGetLastStateOfCharge(t *testing.T) {
	router := SetupRouter()

	// test with empty LastQPGSResponse (like pre first request)
	LastQPGSResponse = (*messages.QPGSnResponse)(nil)

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "/last/soc", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "{\"BatteryStateOfCharge\":\"null\"}", w.Body.String())

	// test with relistic response
	input := "(1 92932004102443 B 00 237.0 50.01 000.0 00.00 0483 0387 009 51.1 000 069 020.4 000 00942 00792 007 00000010 1 1 060 080 10 00.0 006\xf2\x2d\r"
	actual, err := messages.NewQPGSnResponse(input, 1)
	assert.Equal(t, err, nil)
	LastQPGSResponse = actual

	w = httptest.NewRecorder()
	req, err = http.NewRequest(http.MethodGet, "/last/soc", nil)
	assert.Equal(t, err, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "{\"BatteryStateOfCharge\":\"069\"}", w.Body.String())
}

func TestGetMessage(t *testing.T) {
	router := SetupRouter()

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

	// test first get
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/queue/%s", qidUUID2), nil)
	router.ServeHTTP(w, req)
	want := messages.Message{ID: qidUUID2, Command: "QID2", Payload: ""}
	err := json.Unmarshal(w.Body.Bytes(), &actualBody)
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
	router := SetupRouter()

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
	router := SetupRouter()

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	Queue = []messages.Message{
		{ID: qidUUID1, Command: "QID1", Payload: ""},
		{ID: qidUUID2, Command: "QID1", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}

	// test first deletion
	want := []messages.Message{
		{ID: qidUUID1, Command: "QID1", Payload: ""},
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID2), nil)
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

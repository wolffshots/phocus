package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid" // for generating UUIDs for commands
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"wolffshots/phocus/messages"
)

func TestPostMessage(t *testing.T) {
	router := setupRouter()

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	queue = make([]messages.Message, 0)

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
	assert.Equal(t, want, queue)

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
	assert.Equal(t, want, queue)

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
	assert.Equal(t, want, queue)

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
	assert.Equal(t, want, queue)
}

func TestGetQueue(t *testing.T) {
	router := setupRouter()

	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()

	queue = []messages.Message{
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
	assert.Equal(t, want, queue)
}

func TestGetHealth(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "UP", w.Body.String())
}

func TestGetMessage(t *testing.T) {
	router := setupRouter()

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	qidUUID4 := uuid.New()

	queue = []messages.Message{
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
	router := setupRouter()

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	queue = []messages.Message{
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
	assert.Equal(t, want, queue)

	// test empty case
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req) // ensure no deadlock
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, want, queue)
}

func TestDeleteMessage(t *testing.T) {
	router := setupRouter()

	qidUUID1 := uuid.New()
	qidUUID2 := uuid.New()
	qidUUID3 := uuid.New()
	queue = []messages.Message{
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
	assert.Equal(t, want, queue)

	// test second deletion
	want = []messages.Message{
		{ID: qidUUID3, Command: "QID3", Payload: ""},
	}
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID1), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, want, queue)

	// test third deletion
	want = make([]messages.Message, 0)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID3), nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, want, queue)

	// test invalid deletion
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodDelete, fmt.Sprintf("/queue/%s", qidUUID3), nil)
	router.ServeHTTP(w, req) // ensure no deadlock and test negative case
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, want, queue) // assert that it hasn't changed
}

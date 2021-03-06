package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPairDeviec struct{}

func (mockPairDeviec) Pair(p Pair) error {
	return nil
}

func TestPairDeviceHandler(t *testing.T) {

	payload := new(bytes.Buffer)
	json.NewEncoder(payload).Encode(Pair{
		DeviceID: 1234,
		UserID:   4321,
	})
	req := httptest.NewRequest(http.MethodPost, "/pair-device", payload)
	rec := httptest.NewRecorder()

	handler := CustomHandlerFunc(PairDeviceHandler(mockPairDeviec{}))
	handler.ServeHTTP(rec, req)

	if http.StatusOK != rec.Code {
		t.Error("Expect 200 OK but got ", rec.Code)
	}

	expected := fmt.Sprintf("%s\n", `{"status":"active"}`)
	if rec.Body.String() != expected {
		t.Errorf("expected %q but got %q\n", expected, rec.Body.String())
	}
}

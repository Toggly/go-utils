package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestRequestIdGenerateIfEmpty(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := w.Header().Get(XRequestID)
		assert.NotEmpty(t, val, "Request id is not empty")
	})

	handler := RequestIDCtx(nextHandler)
	handler.ServeHTTP(rr, req)
}

func TestRequestIdGenerateIfSet(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// generate request id
	reqID := "UNIT-REQ-" + uuid.NewV4().String()
	req.Header.Add(XRequestID, reqID)

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := w.Header().Get(XRequestID)
		assert.Equal(t, reqID, val, "Request ids are equal")
	})

	handler := RequestIDCtx(nextHandler)
	handler.ServeHTTP(rr, req)
}
func TestVersionContext(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		val := ctx.Value(ContextVersionKey).(string)
		assert.Equal(t, "v1.0", val, "Versions are equal")
	})

	handler := VersionCtx("v1.0")(nextHandler)
	handler.ServeHTTP(rr, req)
}
func TestServiceInfoInHeaders(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceName := w.Header().Get(XServiceName)
		serviceVersion := w.Header().Get(XServiceVersion)
		assert.Equal(t, "shield", serviceName, "Names are equal")
		assert.Equal(t, "1.0.0", serviceVersion, "Versions are equal")
	})

	handler := ServiceInfo("shield", "1.0.0")(nextHandler)
	handler.ServeHTTP(rr, req)
}

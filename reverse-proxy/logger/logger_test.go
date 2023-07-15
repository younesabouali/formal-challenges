package logger

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/younesabouali/formal-challenges/reverse-proxy/Interceptor"
)

func TestInspector(t *testing.T) {
	testCases := []struct {
		desc       string
		method     string
		body       []byte
		config     []Interceptor.InterceptRequest
		wantStatus bool
	}{
		{
			desc:   "want to log",
			method: http.MethodPost,
			config: []Interceptor.InterceptRequest{
				{
					Method: []string{http.MethodPost},
				},
			},
			wantStatus: true,
		},
		{
			desc:   "Don't want to log",
			method: http.MethodGet,
			config: []Interceptor.InterceptRequest{
				{
					Method: []string{http.MethodPost},
				},
			},
			wantStatus: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest(tC.method, "http://example.com/foo", bytes.NewBuffer(tC.body))
			loggerInstance := &Logger{
				FlagLogRequest:  tC.config,
				FlatLogResponse: tC.config,
			}
			type email struct{ Email string }
			emailResponse, _ := json.Marshal(email{})

			w := httptest.NewRecorder()

			w.Header().Set("Content-Type", "application/json")
			w.Write(emailResponse)
			loggedRequest := loggerInstance.LogRequest(req)
			loggedResponse := loggerInstance.LogResponse(req, w.Result(), string(emailResponse))

			if loggedRequest != tC.wantStatus || loggedResponse != tC.wantStatus {
				t.Errorf("Expected status code %v, got %v %v", tC.wantStatus, loggedRequest, loggedResponse)
			}
		})
	}
}

package inspector

import (
	"bytes"
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
			desc:   "want to inspect",
			method: http.MethodPost,
			config: []Interceptor.InterceptRequest{
				{
					Method: []string{http.MethodPost},
				},
			},
			wantStatus: true,
		},
		{
			desc:   "Don't want to inspect",
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
			inspectorInstance := &Inspector{
				Config: tC.config,
			}

			resp := inspectorInstance.Inspect(req)

			if resp != tC.wantStatus {
				t.Errorf("Expected status code %v, got %v", tC.wantStatus, resp)
			}
		})
	}
}

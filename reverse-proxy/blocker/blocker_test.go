package blocker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/younesabouali/formal-challenges/reverse-proxy/Interceptor"
)

func TestBlocker(t *testing.T) {
	testCases := []struct {
		desc       string
		method     string
		body       []byte
		config     []Interceptor.InterceptRequest
		wantStatus int
	}{
		{
			desc:   "Blocked POST request",
			method: http.MethodPost,
			config: []Interceptor.InterceptRequest{
				{
					Method: []string{http.MethodPost},
				},
			},
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			desc:   "Allowed GET request",
			method: http.MethodGet,
			config: []Interceptor.InterceptRequest{
				{
					Method: []string{http.MethodPost},
				},
			},
			wantStatus: http.StatusOK,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			req := httptest.NewRequest(tC.method, "http://example.com/foo", bytes.NewBuffer(tC.body))
			w := httptest.NewRecorder()
			block := &Blocker{
				Config: tC.config,
			}

			block.Block(w, req)
			resp := w.Result()

			if resp.StatusCode != tC.wantStatus {
				t.Errorf("Expected status code %d, got %d", tC.wantStatus, resp.StatusCode)
			}
		})
	}
}

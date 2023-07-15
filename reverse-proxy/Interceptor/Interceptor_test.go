package Interceptor

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockRequest(method, url, body string, headers map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBufferString(body))

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req
}

func TestMethodFlaging(t *testing.T) {
	testCases := []struct {
		desc       string
		method     string
		body       []byte
		config     []InterceptRequest
		wantResult bool
	}{
		{
			desc:   "In case method is empty flag everything",
			method: http.MethodPost,
			config: []InterceptRequest{
				{
					Method: []string{},
				},
			},
			wantResult: true,
		},
		{
			desc:   "Method not in the flag list",
			method: http.MethodGet,
			config: []InterceptRequest{
				{
					Method: []string{http.MethodPost},
				},
			},
			wantResult: false,
		},

		{
			desc:   "Method in the flag list",
			method: http.MethodPost,
			config: []InterceptRequest{
				{
					Method: []string{http.MethodPost},
				},
			},
			wantResult: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req := httptest.NewRequest(tC.method, "http://example.com/foo", bytes.NewBuffer(tC.body))
			result := isMethodFlagged(req, tC.config[0])

			if result != tC.wantResult {
				t.Errorf("Expected the flag to be %v, got %v", tC.wantResult, result)
			}

		})

	}
}

func TestHeaderFlaging(t *testing.T) {
	testCases := []struct {
		desc       string
		header     string
		body       []byte
		config     []InterceptRequest
		wantResult bool
	}{
		{
			desc: "In case empty is empty flag everything",

			header: "header-2",
			config: []InterceptRequest{
				{
					Headers: InterceptionCondition{},
					// Method: []string{},
				},
			},
			wantResult: true,
		},
		{
			desc:   "Header not in the flag list",
			header: "header-2",
			config: []InterceptRequest{
				{
					Headers: InterceptionCondition{
						Includes: []string{"header-1"},
					},
					// Method: []string{},
				},
			},
			wantResult: false,
		},

		{
			desc:   "Header in the flag list",
			header: "header-1",
			config: []InterceptRequest{
				{
					Headers: InterceptionCondition{
						Includes: []string{"header-1"},
					},
					// Method: []string{},
				},
			},
			wantResult: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", bytes.NewBuffer(tC.body))
			// tC.config[0].
			req.Header.Add(tC.header, "Header Value")

			result := isHeaderFlagged(req, tC.config[0])

			if result != tC.wantResult {
				t.Errorf("Expected the flag to be %v, got %v", tC.wantResult, result)
			}

		})

	}
}

func TestQueryFlaging(t *testing.T) {
	testCases := []struct {
		desc       string
		query      string
		body       []byte
		config     []InterceptRequest
		wantResult bool
	}{
		{
			desc: "In case query is empty flag everything",

			query: "header-2",
			config: []InterceptRequest{
				{
					Query: InterceptionCondition{},
					// Method: []string{},
				},
			},
			wantResult: true,
		},
		{
			desc:  "Query not in the flag list",
			query: "header-2",
			config: []InterceptRequest{
				{
					Query: InterceptionCondition{
						Includes: []string{"header-1"},
					},
					// Method: []string{},
				},
			},
			wantResult: false,
		},

		{
			desc:  "Query in the flag list",
			query: "header-1",
			config: []InterceptRequest{
				{
					Query: InterceptionCondition{
						Includes: []string{"header-1"},
					},
				},
			},
			wantResult: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			appendedQuery := ""
			if tC.query != "" {
				appendedQuery = "?" + tC.query + "=hello"
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com/foo"+appendedQuery, bytes.NewBuffer(tC.body))
			// tC.config[0].

			result := isQueryParamsFlagged(req, tC.config[0])

			if result != tC.wantResult {
				t.Errorf("Expected the flag to be %v, got %v", tC.wantResult, result)
			}

		})

	}
}

func TestPathFlaging(t *testing.T) {
	testCases := []struct {
		desc       string
		path       string
		body       []byte
		config     []InterceptRequest
		wantResult bool
	}{
		{
			desc: "In case path is empty flag everything",

			path: "/header-2",
			config: []InterceptRequest{
				{
					Path: InterceptionCondition{},
					// Method: []string{},
				},
			},
			wantResult: true,
		},
		{
			desc: "path not in the flag list",
			path: "/header-2",
			config: []InterceptRequest{
				{
					Path: InterceptionCondition{
						Includes: []string{"header-1"},
					},
					// Method: []string{},
				},
			},
			wantResult: false,
		},

		{
			desc: "path in the flag list",
			path: "/header-1",
			config: []InterceptRequest{
				{
					Path: InterceptionCondition{
						Includes: []string{"header-1"},
					},
				},
			},
			wantResult: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			appendedPath := ""
			if tC.path != "" {
				appendedPath = tC.path
			}

			req := httptest.NewRequest(http.MethodGet, "http://example.com"+appendedPath, bytes.NewBuffer(tC.body))
			// tC.config[0].

			result := isPathFlagged(req, tC.config[0])

			if result != tC.wantResult {
				t.Errorf("Expected the flag to be %v, got %v", tC.wantResult, result)
			}

		})

	}
}

func TestBodyFlaging(t *testing.T) {
	type email struct{ Email string }
	emailBody, _ := json.Marshal(email{})
	testCases := []struct {
		desc       string
		body       []byte
		config     []InterceptRequest
		wantResult bool
	}{
		{
			desc: "In case body is empty flag everything",

			config:     []InterceptRequest{{}},
			body:       emailBody,
			wantResult: true,
		},
		{
			desc: "body not in the flag list",
			body: emailBody,
			config: []InterceptRequest{
				{
					Body: InterceptionCondition{
						Includes: []string{"Name"},
					},
					// Method: []string{},
				},
			},
			wantResult: false,
		},

		{
			desc: "body in the flag list",
			body: emailBody,
			config: []InterceptRequest{
				{
					Body: InterceptionCondition{
						Includes: []string{"Email"},
					},
				},
			},
			wantResult: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "http://example.com", bytes.NewBuffer(tC.body))
			req.Header.Set("Content-Type", "application/json")
			// tC.config[0].

			result := isBodyFlagged(req, tC.config[0])

			if result != tC.wantResult {
				t.Errorf("Expected the flag to be %v, got %v", tC.wantResult, result)
			}

		})

	}
}

func TestFlagingList(t *testing.T) {
	type email struct{ Email string }
	emailBody, _ := json.Marshal(email{})
	config := []InterceptRequest{{
		Method: []string{http.MethodGet},
		Body:   InterceptionCondition{Includes: []string{"Name"}},
	}, {

		Method: []string{http.MethodPost},
		Path:   InterceptionCondition{Includes: []string{"name"}},
	}}
	testCases := []struct {
		desc       string
		body       []byte
		method     string
		path       string
		config     []InterceptRequest
		wantResult bool
	}{
		{
			desc: "testing a request that is not flagged",

			config: config,
			body:   emailBody,
			method: http.MethodPost,

			wantResult: false,
		},
		{
			desc:   "testing a request that is flagged",
			body:   emailBody,
			config: config,

			method:     http.MethodPost,
			path:       "/name",
			wantResult: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			appendedPath := ""
			if tC.path != "" {
				appendedPath += tC.path
			}

			req := httptest.NewRequest(tC.method, "http://example.com"+appendedPath, bytes.NewBuffer(tC.body))
			req.Header.Set("Content-Type", "application/json")
			// tC.config[0].

			result := IsRequestFlagged(req, tC.config)

			if result != tC.wantResult {
				t.Errorf("Expected the flag to be %v, got %v", tC.wantResult, result)
			}

		})

	}
}

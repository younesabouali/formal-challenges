package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

type message struct {
	Password string
}

func TestReverseProxy(t *testing.T) {
	// Set up a test server that will receive the redirected requests
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write a response so we can verify the request was proxied
		m := message{Password: "Heelo"}
		json.NewEncoder(w).Encode(m)
	}))
	defer targetServer.Close()

	// Use your `NewReverseProxy` function to set up the proxy
	// Note: You will need to replace `proxyInstance.NewReverseProxy` with
	// the actual name of your reverse proxy creation function
	config := getConfig()
	proxyInstance := Proxy{config: &config}
	go startTargetServer(config.TargetUrl)
	var wg sync.WaitGroup
	proxyServer := proxyInstance.NewReverseProxy(&wg)
	wg.Add(1)
	// Set up a second server that uses the reverse proxy
	defer proxyServer.Close()

	// Make a request to the proxy server
	res, err := http.Get("http://localhost:" + config.Port)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Read the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Unmarshal the response
	var m message
	err = json.Unmarshal(body, &m)
	if err != nil {
		t.Fatal(err)
	}
	val := reflect.ValueOf(m)
	// Check if the response came from the target server
	for _, censoredWord := range config.CensorFields {
		if fieldVal := val.FieldByName(censoredWord); fieldVal.IsValid() && fieldVal.Interface() != "********" {
			t.Errorf("Expected the '%s' to be '********', got: %s", censoredWord, val)
		}
	}
}

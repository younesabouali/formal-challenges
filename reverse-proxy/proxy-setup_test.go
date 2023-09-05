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
	Port := "3000"
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := message{Password: "Heelo"}
		json.NewEncoder(w).Encode(m)
	}))
	defer targetServer.Close()

	config := getConfig()
	config.Port = Port
    config.TargetUrl ="http://localhost:3001"
	proxyInstance := Proxy{config: &config}
	go startTargetServer("http://localhost:3001")
	var wg sync.WaitGroup
	proxyServer := proxyInstance.NewReverseProxy(&wg)
	wg.Add(1)
	defer proxyServer.Close()

	res, err := http.Get("http://localhost:" + Port)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	// Read the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var m message
	err = json.Unmarshal(body, &m)
	if err != nil {
		t.Fatal(err)
	}
	val := reflect.ValueOf(m)
	for _, censoredWord := range config.CensorFields {
		if fieldVal := val.FieldByName(censoredWord); fieldVal.IsValid() && fieldVal.Interface() != "********" {
			t.Errorf("Expected the '%s' to be '********', got: %s", censoredWord, val)
		}
	}
}

package main

import (
	"sync"
)

func main() {

	// Read the config file
	// todo: accept file path from the args
	proxyConfig := getConfig()
	if proxyConfig.SpinTarget == true {
		go startTargetServer(proxyConfig.TargetUrl)
	}
	// Create a new reverse proxy
	var wg sync.WaitGroup
	wg.Add(1)
	proxyInstance := Proxy{config: &proxyConfig}
	go proxyInstance.NewReverseProxy(&wg)
	wg.Wait()
	// log.Fatal(server.ListenAndServe())
}

// NewReverseProxy creates a new reverse proxy handler

package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/younesabouali/formal-challenges/reverse-proxy/Interceptor"
)

type ProxyConfig struct {
	Port       string
	TargetUrl  string
	SpinTarget bool

	CensorFields []string

	InspectRequest []Interceptor.InterceptRequest
	LogRequest     []Interceptor.InterceptRequest
	LogResponse    []Interceptor.InterceptRequest
	BlockRequest   []Interceptor.InterceptRequest
}

func getConfig() ProxyConfig {

	filePath := "proxy-config.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var proxyConfig ProxyConfig
	err = json.Unmarshal(data, &proxyConfig)
	if err != nil {

		log.Fatal("Error parsing proxy-config ", err)
	}
	return proxyConfig
}

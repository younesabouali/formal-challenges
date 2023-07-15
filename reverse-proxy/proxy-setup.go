package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/younesabouali/formal-challenges/reverse-proxy/blocker"
	"github.com/younesabouali/formal-challenges/reverse-proxy/censor"
	"github.com/younesabouali/formal-challenges/reverse-proxy/inspector"
	"github.com/younesabouali/formal-challenges/reverse-proxy/logger"
)

type Proxy struct {
	config *ProxyConfig
}

func (P *Proxy) NewReverseProxy(wg *sync.WaitGroup) *http.Server {
	// Parse the target URL
	targetURL, _ := url.Parse(P.config.TargetUrl)

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	http.HandleFunc("/", P.handler(proxy))
	log.Printf("Reverse proxy server is running on Port: %v and redirecting to %v", P.config.Port, P.config.TargetUrl)
	srv := &http.Server{Addr: ":" + P.config.Port}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
		wg.Done()
	}()

	return srv
}
func (P *Proxy) handler(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		p.ModifyResponse = P.modifyResponse(r)
		loggerInstance := logger.Logger{FlagLogRequest: P.config.LogRequest}
		loggerInstance.LogRequest(r)

		blockerInstance := blocker.Blocker{Config: P.config.BlockRequest}
		if blockerInstance.Block(w, r) {
			return
		}

		inspectorInstance := inspector.Inspector{Config: P.config.InspectRequest}
		inspectorInstance.Inspect(r)
		p.ServeHTTP(w, r)
	}
}

func (P *Proxy) modifyResponse(r *http.Request) func(res *http.Response) error {
	return func(res *http.Response) error {

		if res.ContentLength == 0 {
			return nil
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = res.Body.Close()
		if err != nil {
			return err
		}

		mBody := censor.Censor(P.config.CensorFields)(string(body))
		res.Body = ioutil.NopCloser(strings.NewReader(mBody))

		loggerInstance := logger.Logger{FlatLogResponse: P.config.LogResponse}
		loggerInstance.LogResponse(r, res, mBody)
		res.ContentLength = int64(len(mBody))
		res.Header.Set("Content-Length", strconv.Itoa(len(mBody)))

		return nil
	}
}

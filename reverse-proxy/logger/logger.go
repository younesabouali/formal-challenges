package logger

import (
	"log"
	"net/http"

	"github.com/younesabouali/formal-challenges/reverse-proxy/Interceptor"
)

type Logger struct {
	FlagLogRequest  []Interceptor.InterceptRequest
	FlatLogResponse []Interceptor.InterceptRequest
}

func (l *Logger) LogRequest(r *http.Request) bool {
	log.Printf("Received request %s \n  %s \n %s\n Headers: %v\n", r.Method, r.Host, r.RemoteAddr, r.Header)
	if value := Interceptor.IsRequestFlagged(r, l.FlagLogRequest); value {
		return true
	}
	return false
}
func (l *Logger) LogResponse(r *http.Request, res *http.Response, mBody string) bool {

	if value := Interceptor.IsRequestFlagged(r, l.FlatLogResponse); value {
		return true
	}
	log.Printf("Received response %s\nHeaders: %v\n", res.Status, res.Header)
	log.Printf("Body: %s\n", string(mBody))
	return false
}

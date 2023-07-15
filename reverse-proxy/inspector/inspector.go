package inspector

import (
	"log"
	"net/http"

	"github.com/younesabouali/formal-challenges/reverse-proxy/Interceptor"
)

type Inspector struct {
	Config []Interceptor.InterceptRequest
}

func (i *Inspector) Inspect(r *http.Request) bool {
	if value := Interceptor.IsRequestFlagged(r, i.Config); value {
		log.Printf("Inspection of GET request, URL: %s\n", r.URL.String())
		return true
		// You can inspect the request headers, body, etc. here
	}
	return false
}

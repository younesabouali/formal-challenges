package blocker

import (
	"fmt"
	"net/http"

	"github.com/younesabouali/formal-challenges/reverse-proxy/Interceptor"
)

type Blocker struct {
	Config []Interceptor.InterceptRequest
}

func (b *Blocker) Block(w http.ResponseWriter, r *http.Request) bool {
	if value := Interceptor.IsRequestFlagged(r, b.Config); value {
		fmt.Print("Request Flagged")
		// We're blocking all POST requests
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return true
	}
	return false

}

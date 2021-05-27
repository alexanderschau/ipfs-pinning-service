package pinning

import (
	"fmt"
	"net/http"
)

// Handler - HTTP Handler Function
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	w.Write([]byte("hello world"))
}

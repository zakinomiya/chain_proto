package gateway

import (
	"fmt"
	"log"
	"net/http"
)

type HTTPServer struct {
	mux *http.ServeMux
}

func NewHTTPServer() *HTTPServer {
	h := &HTTPServer{
		mux: &http.ServeMux{},
	}

	h.mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "world")
	})

	return h
}

func (h *HTTPServer) Start(port string) {
	if err := http.ListenAndServe(":"+port, h.mux); err != nil {
		log.Fatalln("error: failed to start http server")
	}
}

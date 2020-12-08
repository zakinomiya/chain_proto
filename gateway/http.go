package gateway

import (
	"fmt"
	"log"
	"net/http"
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer(port string) *HTTPServer {
	mux := &http.ServeMux{}

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "world")
	})

	return &HTTPServer{
		&http.Server{
			Addr:    ":" + port,
			Handler: mux,
		},
	}
}

func (h *HTTPServer) Start() {
	if err := h.ListenAndServe(); err != nil {
		log.Fatalln("error: failed to start http server")
	}
}

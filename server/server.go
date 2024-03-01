package server

import (
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
}

func RunAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	helloHandler := func(w http.ResponseWriter, req *http.Request) {
		_, err := io.WriteString(w, "Hello, world!\n")
		if err != nil {
			log.Printf("Cannot write hello world!. %v", err)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/", helloHandler)

	log.Println("JSON Api server running on port: ", s.listenAddr)
	err := http.ListenAndServe(s.listenAddr, r)
	if err != nil {
		log.Fatal(err)
	}
}

package server

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/server/router"

	"gorm.io/gorm"
)

type APIServer struct {
	listenAddr string
	database   *gorm.DB
	router     *mux.Router
}

func NewAPIServer(listenAddr string, db *gorm.DB) *APIServer {
	apiRouter := mux.NewRouter()

	router.HandleCategoryRouter(db, apiRouter)

	return &APIServer{
		listenAddr: listenAddr,
		database:   db,
		router:     apiRouter,
	}
}

func (s *APIServer) Run() {
	var wg sync.WaitGroup
	wg.Add(1)
	go s.startServer(&wg)
	wg.Wait()
}

func (s *APIServer) startServer(wg *sync.WaitGroup) {
	defer wg.Done()

	srv := &http.Server{
		Handler:      s.router,
		Addr:         s.listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("RAT api server running on port:", s.listenAddr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", s.listenAddr, err)
	}
}

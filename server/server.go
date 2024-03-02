package server

import (
	"errors"
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
	appRouter := mux.NewRouter()

	api := appRouter.PathPrefix("/api").Subrouter()

	handlersRouter := []func(*gorm.DB, *mux.Router, *sync.WaitGroup){
		router.HandleCategoryRouter,
		router.HandleCustomFieldRouter,
	}
	wg := sync.WaitGroup{}
	wg.Add(len(handlersRouter))

	for _, handler := range handlersRouter {
		go handler(db, api, &wg)
	}

	wg.Wait()

	return &APIServer{
		listenAddr: listenAddr,
		database:   db,
		router:     appRouter,
	}
}

func (s *APIServer) Run() {
	s.startServer()
}

func (s *APIServer) startServer() {
	srv := &http.Server{
		Handler:      s.router,
		Addr:         s.listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("RAT api server running on port:", s.listenAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not listen on %s: %v\n", s.listenAddr, err)
	}
}

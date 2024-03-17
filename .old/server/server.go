package server

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/kwiats/rate-all-things/pkg/config"
	"github.com/kwiats/rate-all-things/server/middleware"

	"github.com/gorilla/mux"
	"github.com/kwiats/rate-all-things/server/router"

	"gorm.io/gorm"
)

type APIServer struct {
	listenAddr string
	database   *gorm.DB
	router     *mux.Router
	config     *config.Config
}

func NewAPIServer(listenAddr string, db *gorm.DB, config *config.Config) *APIServer {
	appRouter := mux.NewRouter()

	api := appRouter.PathPrefix("/api").Subrouter()

	api.Use(middleware.LoggingMiddleware)

	handlersRouter := []func(*gorm.DB, *mux.Router, *sync.WaitGroup){
		router.HandleCategoryRouter,
		router.HandleCustomFieldRouter,
		router.HandleAuthRouter,
		router.HandleUserRouter,
		router.HandleReviewItemRouter,
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
		config:     config,
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

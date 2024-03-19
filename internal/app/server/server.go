package server

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
	"tit/internal/app/middleware"
	"tit/internal/config"
	"tit/internal/router"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Server struct {
	listenAddr string
	database   *gorm.DB
	router     *mux.Router
	config     *config.Config
}

func NewServer(listenAddr string, db *gorm.DB, config *config.Config) *Server {
	appRouter := mux.NewRouter()
	appRouter.Use(middleware.LoggingMiddleware)

	handlersRouter := []func(*gorm.DB, *mux.Router, *sync.WaitGroup){
		router.HandleUserRouter,
		router.HandleWebSocket,
		router.HandleAuthRouter,
	}
	wg := sync.WaitGroup{}
	wg.Add(len(handlersRouter))

	for _, handler := range handlersRouter {
		go handler(db, appRouter, &wg)
	}

	wg.Wait()

	return &Server{
		listenAddr: listenAddr,
		database:   db,
		router:     appRouter,
		config:     config,
	}
}

func (s *Server) Run() {
	s.startServer()
}

func (s *Server) startServer() {
	srv := &http.Server{
		Handler:      s.router,
		Addr:         s.listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting server on http://localhost%s/", s.listenAddr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Could not listen on %s: %v\n", s.listenAddr, err)
	}
}

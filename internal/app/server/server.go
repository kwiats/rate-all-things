package server

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
	"tit/internal/app/middleware"
	"tit/internal/auth"
	"tit/internal/chat"
	"tit/internal/config"
	"tit/internal/user"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

type Server struct {
	listenAddr string
	database   *gorm.DB
	router     http.Handler
	config     *config.Config
}

func NewServer(listenAddr string, db *gorm.DB, config *config.Config) *Server {
	appRouter := mux.NewRouter()
	appRouter.Use(middleware.LoggingMiddleware)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{config.Envs.AllowedDomains},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	})
	handlersRouter := []func(*gorm.DB, *mux.Router, *sync.WaitGroup){
		user.HandleUserRouter,
		chat.HandleWebSocketRouter,
		chat.HandleChatRouter,
		auth.HandleAuthRouter,
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
		router:     c.Handler(appRouter),
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

package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suryansh74/simplebank/db"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// accounts
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	// transfers
	router.POST("/transfers", server.createTransfer)

	// users
	router.POST("/users", server.createUser)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) StartWithShutdown(ctx context.Context, address string) error {
	srv := &http.Server{
		Addr:    address,
		Handler: server.router,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(context.Background())
	}()

	return srv.ListenAndServe()
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

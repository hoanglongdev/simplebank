package api

import (
	"fmt"
	db "simple_bank/db/sqlc"
	"simple_bank/token"
	"simple_bank/util"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Server struct {
	config util.Config
	maker  token.Maker
	store  db.Store
	router *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config: config,
		maker:  tokenMaker,
		store:  store,
	}

	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := validate.RegisterValidation("currency", validCurrency)
		if err != nil {
			return nil, err
		}
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {

	// add routes to default router
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginuser)

	// add routes to auth router
	authRouters := router.Group("/").Use(authMiddleware(server.maker))
	authRouters.POST("/accounts", server.createAccount)
	authRouters.GET("/accounts/:id", server.getAccount)
	authRouters.GET("/accounts", server.listAccount)
	authRouters.POST("/transfers", server.createTransfer)

	// add router to server
	server.router = router
}

func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResp(err error) gin.H {
	return gin.H{"error": err.Error()}
}

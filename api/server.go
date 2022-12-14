package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/kvn-media/backend_test/db/sqlc"
	docs "github.com/kvn-media/backend_test/docs"
	"github.com/kvn-media/backend_test/token"
	"github.com/kvn-media/backend_test/util"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Simple Backend Bank API
// @version 1.0
// @description A simple bank service.

// @contact.name API Support
// @contact.url https://github.com/kvn-media/backend_test/issues
// @contact.email kevin.subagio@gmail.com

// @host localhost:8080
// @BasePath /

// @securityDefinitions.basic BasicAuth

// @securityDefinitions.apikey authorization
// @in header
// @name Authorization

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}
	//binding custom validator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//registor validator to gin
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	docs.SwaggerInfo.BasePath = "/"
	router := gin.Default()
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.CreateAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.CreateTransfer)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//add routes to router
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

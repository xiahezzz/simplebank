package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/xiahezzz/simplebank/db/sqlc"
	"github.com/xiahezzz/simplebank/db/util"
	"github.com/xiahezzz/simplebank/token"
)

type Server struct {
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
	//帮助我们将请求发送到的正确的处理程序
	router *gin.Engine
}

//创建一个新的Http Server并且建立路由
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot make tokenMaker:%w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	//add routes to router
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)
	authRoutes.POST("/accounts/update", server.updateAccount)
	authRoutes.POST("/accounts/delete", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
}

//服务器启动
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

//错误处理
func errorRequest(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

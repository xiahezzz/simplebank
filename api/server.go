package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/xiahezzz/simplebank/db/sqlc"
)

type Server struct {
	store db.Store
	//帮助我们将请求发送到的正确的处理程序
	router *gin.Engine
}

//创建一个新的HPTTP Server并且建立路由
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	//注册验证器
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	//add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.POST("/accounts/update", server.updateAccount)
	router.POST("/accounts/delete", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
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

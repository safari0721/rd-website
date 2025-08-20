package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/rd-website/pkg/handlers"
)

func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/signup", authHandler.Signup)
	router.POST("/login", authHandler.Login)

	return router
}

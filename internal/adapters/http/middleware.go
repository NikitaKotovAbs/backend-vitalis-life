package http

import(
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"backend/config"
	"time"
)

func CORSNew(config config.CORSConfig) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     config.AllowOrigins,
		AllowMethods:     config.AllowMethods,
		AllowHeaders:     config.AllowHeaders,
		ExposeHeaders:    config.ExposeHeaders,
		AllowCredentials: config.AllowCredentials,
		MaxAge:           time.Second * time.Duration(config.MaxAge),
	})
}
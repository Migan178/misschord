package routers

import "github.com/gin-gonic/gin"

func SetupRouter(app *gin.Engine) {
	v1 := app.Group("/v1")

	setupUsers(v1)
}

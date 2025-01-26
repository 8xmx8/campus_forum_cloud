package svc

import "github.com/gin-gonic/gin"

func (s *ServiceContext) Ping(c *gin.Context) {
	g := Gin{C: c}
	g.Response(OK, "", string("pong"))
}

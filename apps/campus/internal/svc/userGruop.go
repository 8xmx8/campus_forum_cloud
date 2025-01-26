package svc

import (
	"campus_forum_cloud/apps/campus/internal/model"
	"campus_forum_cloud/apps/campus/internal/types"
	"github.com/gin-gonic/gin"
)

func (s *ServiceContext) Register(c *gin.Context) {
	g := Gin{C: c}
	var req types.RegisterUser
	if err := c.BindJSON(&req); err != nil {
		g.Response(Bind, err.Error(), nil)
	}

	err := s.DAO.Register(&model.User{
		Username: req.UserName,
		Password: req.Password,
	})
	if err != nil {
		g.Response(Err, err.Error(), nil)
	}
	g.Response(OK, "success", nil)
}

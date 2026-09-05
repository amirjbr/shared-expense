package http

import (
	"github.com/gin-gonic/gin"
)

type RouteRegistrar interface {
	RegisterRoutes(group *gin.RouterGroup)
}

type Server struct {
	engine *gin.Engine
	addr   string
}

func NewServer(Addr string) *Server {
	engine := gin.Default()
	return &Server{
		engine: engine,
		addr:   Addr,
	}
}

func (s *Server) Register(prefix string, routers ...RouteRegistrar) {
	group := s.engine.Group(prefix)
	for _, router := range routers {
		router.RegisterRoutes(group)
	}
}

func (s *Server) Run() error {
	return s.engine.Run(s.addr)
}

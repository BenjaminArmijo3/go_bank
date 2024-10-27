package app

import (
	"net/http"

	"github.com/BenjaminArmijo3/bank/internal/config"
	"github.com/BenjaminArmijo3/bank/internal/db/store"
	"github.com/BenjaminArmijo3/bank/internal/pkg/utils/token"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Router *gin.Engine
	Store  *store.Store
	Token  *token.JWTToken
	cfg    *config.Config
}

type RouterGroup interface {
	RegisterRoutes(rg *gin.RouterGroup)
}

func (s *Server) RegisterRoutes(prefix string, group RouterGroup) {
	routerGroup := s.Router.Group(prefix)

	group.RegisterRoutes(routerGroup)
}

func NewServer(store *store.Store, cfg *config.Config, token *token.JWTToken) *Server {

	g := gin.Default()
	g.Use(cors.Default())

	return &Server{
		Router: g, Store: store, cfg: cfg, Token: token,
	}
}

func (s *Server) Start(address string) error {
	s.Router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Welcome"})
	})

	return s.Router.Run(address)
}

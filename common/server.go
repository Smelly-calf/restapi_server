package common

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

const (
	defaultGracefulTimeout = 5 * time.Second
)

type HttpServer interface {
	Start() error // 启动服务
	Stop() error  // 停止服务
}

type Router interface {
	Route(*gin.Engine)
}

type HttpConfig struct {
	Router          Router
	Addr            string
	GracefulTimeout time.Duration
}

type httpServer struct {
	srv             *http.Server
	gracefulTimeout time.Duration
}

func NewHttpServer(cfg HttpConfig) *httpServer {
	// 创建引擎
	gin.SetMode(gin.ReleaseMode)
	e := gin.Default()
	e.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 加载路由
	cfg.Router.Route(e)

	gto := defaultGracefulTimeout
	if cfg.GracefulTimeout != 0 {
		gto = cfg.GracefulTimeout
	}

	return &httpServer{
		srv: &http.Server{
			Addr:    cfg.Addr,
			Handler: e,
		},
		gracefulTimeout: gto,
	}
}

func (s *httpServer) Start() error {
	log.Println("http server starting...")
	return s.srv.ListenAndServe()
}

func (s *httpServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.gracefulTimeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal("fail to stop http server", err)
		return err
	}
	log.Println("http server stopped")
	return nil
}

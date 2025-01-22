package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RouteFactory func(router *gin.Engine)

type RouterOption struct {
	RecoveryDisabled bool
	LoggerForced     bool
	OriginsAllowed   []string
}

func MustServe(endpoint string, factory RouteFactory, option ...RouterOption) {
	if err := Serve(endpoint, factory, option...); err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("Failed to serve API")
	}
}

func Serve(endpoint string, factory RouteFactory, option ...RouterOption) error {
	router := newRouter(factory, option...)

	server := http.Server{
		Addr:    endpoint,
		Handler: router,
	}

	return server.ListenAndServe()
}

func newRouter(factory RouteFactory, option ...RouterOption) *gin.Engine {
	var opt RouterOption
	if len(option) > 0 {
		opt = option[0]
	}

	router := gin.New()

	if !opt.RecoveryDisabled {
		router.Use(gin.Recovery())
	}

	router.Use(newCorsMiddleware(opt.OriginsAllowed))

	if opt.LoggerForced || logrus.IsLevelEnabled(logrus.DebugLevel) {
		router.Use(gin.Logger())
	}

	factory(router)

	return router
}

func newCorsMiddleware(origins []string) gin.HandlerFunc {
	conf := cors.DefaultConfig()
	conf.AllowMethods = append(conf.AllowMethods, "OPTIONS")
	conf.AllowHeaders = append(conf.AllowHeaders, "*")

	if len(origins) == 0 {
		conf.AllowAllOrigins = true
	} else {
		conf.AllowOrigins = origins
	}

	return cors.New(conf)
}

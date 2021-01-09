package http

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Config(r *gin.Engine) {
	// collector apis, compatible with open-falcon
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", ping)
		v1.GET("/pid", pid)
		v1.GET("/addr", addr)
	}

	pprof.Register(r, "/debug/pprof")
}

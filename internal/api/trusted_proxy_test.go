package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestConfigureTrustedProxiesUsesRightmostUntrustedForwardedIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	require.NoError(t, configureTrustedProxies(engine, []string{"172.16.0.0/12"}))
	var clientIP string
	engine.GET("/ip", func(ctx *gin.Context) {
		clientIP = ctx.ClientIP()
		ctx.Status(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodGet, "/ip", nil)
	req.RemoteAddr = "172.18.0.4:12345"
	req.Header.Set("X-Forwarded-For", "198.51.100.99, 203.0.113.8")
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	require.Equal(t, "203.0.113.8", clientIP)
}

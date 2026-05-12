package proxy

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProxyConfig struct {
	AuthServiceURL string
	UserServiceURL string
}

func ProxyRequest(c *gin.Context, targetURL string) {
	method := c.Request.Method
	path := c.Request.URL.RequestURI()

	req, err := http.NewRequest(method, targetURL+path, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to create proxy request"})
		return
	}

	copyHeaders(c.Request.Header, req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	copyHeaders(resp.Header, c.Writer.Header())
	c.Writer.WriteHeader(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

func copyHeaders(src http.Header, dst http.Header) {
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func ForwardToAuthService(authServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ProxyRequest(c, authServiceURL)
	}
}

func ForwardToUserService(userServiceURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ProxyRequest(c, userServiceURL)
	}
}

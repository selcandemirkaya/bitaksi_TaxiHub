package main

//gin
import (
	"bitaksi_taxihub/gateway/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	driverBase := os.Getenv("DRIVER_SERVICE_URL")
	if driverBase == "" {
		driverBase = "http://driver-service:8081" //"http://localhost:8081"
	}
	targetURL, err := url.Parse(driverBase)
	if err != nil {
		log.Fatalf("invalid DRIVER_SERVICE_URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)

		req.Host = targetURL.Host
	}
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		log.Printf("proxy error: %v", e)
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"driver-service unreachable"}`))
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.RateLimitMiddleware(), middleware.JWTMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "gateway ok"})
	})

	r.Any("/drivers", func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
	r.Any("/drivers/*any", func(c *gin.Context) {
		c.Request.URL.Path = strings.Replace(c.Request.URL.Path, "/drivers/", "/drivers/", 1)
		proxy.ServeHTTP(c.Writer, c.Request)
	})

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("gateway running on: %s -> %s", port, driverBase)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("gateway failed: %v", err)
	}
}

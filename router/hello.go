// router/hello.go

package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTP GET request to test the root router
func HelloGo(c *gin.Context) {

	c.JSON(http.StatusCreated, gin.H{"message": "Welcome to my AD management site."})
}

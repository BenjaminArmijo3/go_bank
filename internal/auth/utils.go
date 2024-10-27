package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)



func GetActiveUser(c *gin.Context) (int64, error) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "not authorized to access resources"})
		return 0, fmt.Errorf("error occurred:")
	}

	userID, ok := userId.(int64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "ecnountered an isue"})
		return 0, fmt.Errorf("error occurred:")
	}

	return userID, nil
}

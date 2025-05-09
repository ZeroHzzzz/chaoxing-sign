package user

import (
	"chaoxing/internal/services"

	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Email    string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var (
		req loginReq
		err error
	)

	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	user, err := services.GetUserByEmailPass(c, req.Email, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	if user == nil {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful", "user": user})
}

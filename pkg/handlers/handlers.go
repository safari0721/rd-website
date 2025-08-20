package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/rd-website/pkg/services"
)

type signupRequest struct {
	Name     string `json:"name" binding:"required"`
	AgentID  string `json:"agent_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginRequest struct {
	AgentID  string `json:"agent_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthHandler struct {
	AuthService *services.AuthService
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var input signupRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(input.Name) < 1 || len(input.Name) > 100 || len(input.AgentID) < 3 || len(input.AgentID) > 64 || len(input.Password) < 8 || len(input.Password) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	err := h.AuthService.Signup(input.Name, input.AgentID, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "signup failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Signup successful"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input loginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(input.AgentID) < 3 || len(input.AgentID) > 64 || len(input.Password) < 8 || len(input.Password) > 128 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := h.AuthService.Login(input.AgentID, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Return non-sensitive user info
	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"agent_name": user.Name,
		"agent_id":   user.AgentID,
	})
}

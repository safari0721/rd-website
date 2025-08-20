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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.AuthService.Signup(input.Name, input.AgentID, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Signup failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Signup successful"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input loginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

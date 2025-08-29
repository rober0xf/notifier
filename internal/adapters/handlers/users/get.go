package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *userHandler) GetUser(email string, c *gin.Context) {
	user, err := h.UserService.Get(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching user by email"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *userHandler) GetAllUsers(c *gin.Context) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching all users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (h *userHandler) GetUserByID(id string, c *gin.Context) {
	id_int, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "error processing user id"})
		return
	}

	user, err := h.UserService.GetUserFromID(uint(id_int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching user by id"})
		return
	}
	c.JSON(http.StatusOK, user)
}

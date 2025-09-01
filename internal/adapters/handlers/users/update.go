package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/domain"
)

func (h *userHandler) UpdateUser(c *gin.Context) {
	var input_user struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// parse json
	if err := c.ShouldBindJSON(&input_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request (update user)"})
		return
	}

	userID, err := h.Utils.GetUserIDFromRequest(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// create the user with the new data
	user := &domain.User{
		ID:       userID,
		Name:     input_user.Name,
		Email:    input_user.Email,
		Password: input_user.Password,
	}

	updated_user, err := h.UserService.Update(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// clean to show in response
	updated_user.Password = ""

	c.JSON(http.StatusOK, updated_user)
}

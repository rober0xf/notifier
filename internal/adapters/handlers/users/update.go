package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/adapters/httphelpers"
	"github.com/rober0xf/notifier/internal/domain"
)

func (h *userHandler) Update(c *gin.Context) {
	id_str := c.Param("id")
	if id_str == "" {
		httphelpers.IDParameterNotProvided(c)
		return
	}

	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "error parsing id"})
		return
	}

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

	// userID, err := h.Utils.GetUserIDFromRequest(c.Request)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// create the user with the new data
	user := &domain.User{
		ID:       uint(id),
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

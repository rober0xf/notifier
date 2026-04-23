package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/delivery/http/dto"
	"github.com/rober0xf/notifier/internal/domain/entity"
)

// Me godoc
// @Summary      Get current user
// @Description  Returns the authenticated user's data from the JWT context.
// @Tags         auth
// @Produce      json
// @Success      200  {object}  dto.MeResponse
// @Failure      401  {object}  dto.ErrorResponse  "Unauthorized"
// @Security     BearerAuth
// @Router       /v1/auth/users/me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, dto.MeResponse{
		UserID: userID.(int),
		Email:  email.(string),
		Role:   string(role.(entity.Role)),
	})
}

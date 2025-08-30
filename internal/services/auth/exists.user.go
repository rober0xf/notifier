package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rober0xf/notifier/internal/domain"
)

func (s *Service) ExistsUser(c *gin.Context, email string) (*domain.User, error) {
	return s.AuthRepo.ExistsUser(email)
}

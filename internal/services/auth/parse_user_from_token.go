package auth

import (
	"errors"
	"fmt"

	"github.com/rober0xf/notifier/internal/adapters/httphelpers/dto"
	"github.com/rober0xf/notifier/internal/services/mail"
)

// TODO: fix return
func (s *Service) ParseUserFromToken(token_string string) (*mail.MailSender, error) {
	userID, err := s.ValidateToken(token_string, s.jwtKey)
	if err != nil {
		return nil, err
	}

	_, err = s.UserRepo.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, dto.ErrNotFound) {
			return nil, dto.ErrUserNotFound
		}
		return nil, fmt.Errorf("internal error: %w", err)
	}

	return nil, nil
}

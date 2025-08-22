package payments

import (
	"github.com/rober0xf/notifier/internal/ports"
)

type Service struct {
	Repo ports.PaymentRepository
}

func NewPayments(repo ports.PaymentRepository) *Service {
	return &Service{
		Repo: repo,
	}
}

var _ ports.PaymentService = (*Service)(nil)

package kafka

import (
	"context"
	"events/internal/entities"
)

type PaymentsAdapter struct {
	topic string
}

func (p *PaymentsAdapter) Send(ctx context.Context, payment entities.Payment) {

}
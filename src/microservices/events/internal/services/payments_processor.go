package services

import (
	"fmt"
	"context"

	"events/internal/entities"
)

type PaymentsProcessor struct {

}

func (p *PaymentsProcessor) Process(ctx context.Context, payment entities.Payment) error {
	fmt.Printf("processed Payment %v\n", payment)

	return nil
}
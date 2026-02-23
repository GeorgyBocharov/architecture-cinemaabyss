package services

import (
	"fmt"

	"events/internal/entities"
)

type PaymentsProcessor struct {

}

func (p *PaymentsProcessor) Process(payment entities.Payment) {
	fmt.Printf("processed Payment %v\n", payment)
}
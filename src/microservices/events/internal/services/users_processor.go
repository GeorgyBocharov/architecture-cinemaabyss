package services

import (
	"fmt"
	"context"

	"events/internal/entities"
)

type UsersProcessor struct {

}

func (p *UsersProcessor) Process(ctx context.Context, user entities.User) error {
	fmt.Printf("processed User %v\n", user)

	return nil
}
package services

import (
	"fmt"
	"context"

	"events/internal/entities"
)

type MoviesProcessor struct {

}

func (p *MoviesProcessor) Process(ctx context.Context, movie entities.Movie) error {
	fmt.Printf("processed Movie %v\n", movie)

	return nil
}
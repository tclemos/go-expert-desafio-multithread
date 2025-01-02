package cep

import (
	"context"
	"fmt"
	"time"

	"github.com/tclemos/go-expert-desafio-multithread/pkg/entity"
)

type CepProvider interface {
	Get(context.Context, entity.Cep) (string, error)
}

type CepService struct {
	providers []CepProvider
}

func NewService(providers []CepProvider) *CepService {
	return &CepService{
		providers: providers,
	}
}

func (s *CepService) Get(ctx context.Context, cep entity.Cep, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	respChan := make(chan string)
	for _, provider := range s.providers {
		go func(p CepProvider) {
			resp, err := p.Get(ctx, cep)
			if err == nil {
				respChan <- resp
			}
		}(provider)
	}

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("tempo limite excedido")
	case resp := <-respChan:
		return resp, nil
	}
}

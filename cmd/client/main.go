package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tclemos/go-expert-desafio-multithread/config"
	brasil_api "github.com/tclemos/go-expert-desafio-multithread/internal/infra/brasil-api"
	via_cep "github.com/tclemos/go-expert-desafio-multithread/internal/infra/via-cep"
	"github.com/tclemos/go-expert-desafio-multithread/internal/services/cep"
	"github.com/tclemos/go-expert-desafio-multithread/pkg/entity"
)

func main() {
	ctx := context.Background()
	cfg := config.Config{Timeout: time.Second}

	cepProviders := []cep.CepProvider{
		brasil_api.NewCepProvider(),
		via_cep.NewCepProvider(),
	}
	cepService := cep.NewService(cepProviders)

	fmt.Println("Bem vindo ao buscador de CEPs!")
	fmt.Println("O sistema aceita CEPs apenas no formato 00000-000 ou 00000000")
	for {
		fmt.Println()
		fmt.Println("Digite o CEP que deseja buscar e pressione enter:")

		reader := bufio.NewReader(os.Stdin)
		valorDigitado, _ := reader.ReadString('\n')
		valorDigitado = valorDigitado[:len(valorDigitado)-1]

		cep, err := entity.NewCep(valorDigitado)
		if err != nil {
			fmt.Println(err)
			continue
		}

		response, err := cepService.Get(ctx, *cep, cfg.Timeout)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println()
		fmt.Println(response)
	}
}

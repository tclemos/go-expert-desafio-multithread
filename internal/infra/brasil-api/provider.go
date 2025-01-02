package brasil_api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/tclemos/go-expert-desafio-multithread/pkg/entity"
)

type CepProvider struct{}

func NewCepProvider() *CepProvider {
	return &CepProvider{}
}

func (p *CepProvider) Get(ctx context.Context, cep entity.Cep) (string, error) {
	cepInput := cep.String(false)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://brasilapi.com.br/api/cep/v1/"+cepInput, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	m := make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		return "", err
	}

	if t, found := m["type"]; found && t.(string) == "service_error" {
		return "", fmt.Errorf("CEP não encontrado")
	}

	logradouro := m["street"].(string)
	bairro := m["neighborhood"].(string)
	cidade := m["city"].(string)
	estado := m["state"].(string)

	return fmt.Sprintf("Brasil API: %s, %s, %s - %s, %s", logradouro, bairro, cidade, estado, cep.String(true)), nil
}

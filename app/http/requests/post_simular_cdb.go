package requests

import (
	"errors"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostSimularCdb struct {
    ValorInicial          float64 `json:"valor_inicial" form:"valor_inicial" form:"valor_inicial"`
    ValorAporte          float64 `json:"valor_porte" form:"valor_aporte" form:"valor_aporte"`
    ValorTaxaAnual          float64 `json:"valor_taxa_anual" form: valor_taxa_anual`
    DataInicial           string `json:"data_inicial" form:"data_inicial" form:"data_inicial"`
    DataFinal             string `json:"data_final" form:"data_final" form:"data_final"`
    DiasLiquidezPorAno int `json:"dias_liquidez_por_ano" form:"dias_liquidez_por_ano"`
}
func (r *PostSimularCdb) ValidarData() error {
	_, err := time.Parse("02/01/2006", r.DataInicial) // "02" é o dia, "01" é o mês, "2006" é o ano no Go
	if err != nil {
		//return errors.New("Data inicial, o formato deve ser d/m/Y")
        return errors.New("Erro inexperado")
	}
	_, err = time.Parse("02/01/2006", r.DataFinal) // "02" é o dia, "01" é o mês, "2006" é o ano no Go
	if err != nil {
		//return errors.New("Data final, o formato deve ser d/m/Y")
        return errors.New("Erro inexperado")
	}
    timeInicial, err := time.Parse("02/01/2006", r.DataInicial)
    if err != nil {
        return errors.New("Erro inexperado")
    }
    timeFinal, err := time.Parse("02/01/2006", r.DataFinal)
    if err != nil {
        return errors.New("Erro inexperado!")
    }
    if timeFinal.Before(timeInicial) {
        //return errors.New("Data inicial deve ser menor do que data final!")
        return errors.New("Erro inexperado!")
    }
    limite := 20
    if (timeFinal.Year() - timeInicial.Year()) > limite {
        return errors.New(fmt.Sprintf("Limite de %d anos excedido!", limite))
    }
	return nil
}

func (r *PostSimularCdb) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostSimularCdb) Rules(ctx http.Context) map[string]string {
	return map[string]string{
        "valor_inicial":            "numeric",
        "valor_aporte":            "numeric",
        "data_inicial":             "required|string",
        "data_final":               "required|string",
        "valor_taxa_anual":         "required|numeric",
        "dias_liquidez_por_ano": "required|numeric",
    }
}

func (r *PostSimularCdb) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostSimularCdb) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostSimularCdb) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}

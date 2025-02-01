package requests

import (
	"errors"
	"fmt"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostSimularCdb struct {
    ValorInicial          float64 `json:"valor_inicial" form:"valor_inicial"`
    ValorAporte          float64 `json:"valor_porte" form:"valor_aporte"`
    ValorTaxaAnual          float64 `json:"valor_taxa_anual" form:"valor_taxa_anual"`
    DataInicial           string `json:"data_inicial" form:"data_inicial"`
    DataFinal             string `json:"data_final" form:"data_final"`
    DiasLiquidezPorAno int `json:"dias_liquidez_por_ano" form:"dias_liquidez_por_ano"`
}
func (r *PostSimularCdb) ValidarData() error {
    data_layout := "2006-01-02"
	_, err := time.Parse(data_layout, r.DataInicial) // "02" é o dia, "01" é o mês, "2006" é o ano no Go
	if err != nil {
        fmt.Println(err.Error())
		//return errors.New("Data inicial, o formato deve ser d/m/Y")
        return errors.New("Erro inexperado")
	}
	_, err = time.Parse(data_layout, r.DataFinal) // "02" é o dia, "01" é o mês, "2006" é o ano no Go
	if err != nil {
        fmt.Println(err.Error())
		//return errors.New("Data final, o formato deve ser d/m/Y")
        return errors.New("Erro inexperado")
	}
    timeInicial, err := time.Parse(data_layout, r.DataInicial)
    if err != nil {
        fmt.Println(err.Error())
        return errors.New("Erro inexperado")
    }
    timeFinal, err := time.Parse(data_layout, r.DataFinal)
    if err != nil {
        fmt.Println(err.Error())
        return errors.New("Erro inexperado!")
    }
    if timeFinal.Before(timeInicial) {
        fmt.Println(err.Error())
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
        "valor_taxa_anual":         "required|string",
        "dias_liquidez_por_ano": "required|numeric",
    }
}

func (r *PostSimularCdb) Messages(ctx http.Context) map[string]string {
    return map[string]string{
        "valor_inicial.numeric":             "O valor inicial deve ser numérico",
        "valor_inicial.required":            "O valor inicial é obrigatório",

        "valor_aporte.numeric":              "O valor do aporte deve ser numérico",
        "valor_aporte.required":             "O valor do aporte é obrigatório",

        "data_inicial.required":             "A data inicial é obrigatória",
        "data_inicial.string":               "A data inicial deve ser uma string",

        "data_final.required":               "A data final é obrigatória",
        "data_final.string":                 "A data final deve ser uma string",

        "valor_taxa_anual.required":         "A taxa anual é obrigatória",
        "valor_taxa_anual.string":          "A taxa anual deve ser numérica",

        "dias_liquidez_por_ano.required":    "O número de dias de liquidez por ano é obrigatório",
        "dias_liquidez_por_ano.numeric":     "O número de dias de liquidez por ano deve ser numérico",
    }
}

func (r *PostSimularCdb) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostSimularCdb) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}

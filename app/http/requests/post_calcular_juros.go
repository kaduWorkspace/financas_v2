package requests

import (
	"errors"
	"time"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostCalcularJuros struct {
    ValorInicial          float64 `json:"valor_inicial" form:"valor_inicial"`
    DataInicial           string `json:"data_inicial" form:"data_inicial"`
    DataFinal             string `json:"data_final" form:"data_final"`
    AporteMensal          float64  `json:"aporte_mensal" form:"aporte_mensal"`
    AporteSemestral       float64  `json:"aporte_semestral" form:"aporte_semestral"`
    TipoFrequenciaAumentoAporte  string `json:"tipo_frequencia_aumento_aporte" form:"tipo_frequencia_aumento_aporte"`
    ValorAumentoAporte    float64  `json:"valor_aumento_aporte" form:"valor_aumento_aporte"`
}


func (r *PostCalcularJuros) Authorize(ctx http.Context) error {
	return nil
}
func (r *PostCalcularJuros) ValidarData() error {
	_, err := time.Parse("02/01/2006", r.DataInicial) // "02" é o dia, "01" é o mês, "2006" é o ano no Go

	if err != nil {
		return errors.New("Data inicial, o formato deve ser d/m/Y")
	}
	_, err = time.Parse("02/01/2006", r.DataFinal) // "02" é o dia, "01" é o mês, "2006" é o ano no Go

	if err != nil {
		return errors.New("Data final, o formato deve ser d/m/Y")
	}

    timeInicial, err := time.Parse("02/01/2006", r.DataInicial)
    if err != nil {
        return errors.New("Erro interno")
    }
    timeFinal, err := time.Parse("02/01/2006", r.DataFinal)
    if err != nil {
        return errors.New("Erro interno")
    }
    if timeFinal.Before(timeInicial) {
        return errors.New("Data inicial deve ser menor do que data final!")
    }
    if (timeFinal.Year() - timeInicial.Year()) > 10 {
        return errors.New("Limite de 10 anos excedido")
    }
	return nil
}

func (r *PostCalcularJuros) Rules(ctx http.Context) map[string]string {
	return map[string]string{
        "valor_inicial":            "numeric",
        "data_inicial":             "required|string",
        "data_final":               "required|string",
        "aporte_mensal":            "required|numeric",
        "aporte_semestral":         "numeric",
        "tipo_frequencia_aumento_aporte":   "string|in:semestral,anual",
        "valor_aumento_aporte":     "numeric",
    }
}

func (r *PostCalcularJuros) Messages(ctx http.Context) map[string]string {
    return map[string]string{
        "valor_inicial.numeric":             "Valor inicial deve ser númerico",
        "data_inicial.required":             "Data inicial é obrigatória",
        "data_inicial.string":               "Data inicial deve ser uma string",
        "data_final.required":               "Data final é obrigatória",
        "data_final.string":                 "Data final deve ser uma string",
        "aporte_mensal.required":            "Aporte mensal é obrigatório",
        "aporte_mensal.numeric":             "Aporte mensal deve ser númerico",
        "aporte_semestral.numeric":          "Aporte semestral deve ser númerico",
        "tipo_frequencia_aumento_aporte.string":      "Tipo de frequencia de aumento de aporte deve string",
        "tipo_frequencia_aumento_aporte.in":      "Tipo de frequencia de aumento de aporte deve ser anual ou semestral",
        "valor_aumento_aporte.numeric":              "Valor do aumento do aporte deve ser númerico",
    }
}

func (r *PostCalcularJuros) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostCalcularJuros) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}

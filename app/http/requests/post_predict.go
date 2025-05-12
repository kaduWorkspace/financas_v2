package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type PostPredict struct {
	FutureValue float64 `form:"valor_futuro" json:"valor_futuro"`
    Tax float64 `form:"taxa_juros_anual" json:"taxa_juros_anual"`
    Periods int `form:"periodos" json:"periodos"`
    ContributionOnFirstDay bool `form:"contribuicao_inicio_periodo"`
    InitialValue float64 `form:"valor_inicial" json:"valor_inicial"`
}

func (r *PostPredict) Authorize(ctx http.Context) error {
	return nil
}

func (r *PostPredict) Rules(ctx http.Context) map[string]string {
    return map[string]string{
        "valor_futuro":         "required|string",
        "valor_inicial":         "string",
        "periodos":         "required|string",
        "taxa_juros_anual":         "required|string",
    }
}

func (r *PostPredict) Messages(ctx http.Context) map[string]string {
	return map[string]string{
        "valor_futuro.required": "valor futuro é obrigatório",
        "valor_futuro.numeric": "valor futuro deve ser um número",
        "valor_inicial.numeric": "valor inicial deve ser um número",
        "periodos.required": "periodos é obrigatório",
        "periodos.numeric": "periodos deve ser um número",
        "taxa_juros_anual.required": "taxa juros anual é obrigatório",
        "taxa_juros_anual.numeric": "taxa juros anual deve ser um número",
    }
}

func (r *PostPredict) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *PostPredict) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}

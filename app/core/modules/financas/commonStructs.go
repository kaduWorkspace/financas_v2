package financas

import (
	"fmt"
	"goravel/app/core"
	"math"
	"time"

	"github.com/shopspring/decimal"
)
const decimalPrecision = 16
func DecimalPow(base decimal.Decimal, exponent int64) decimal.Decimal {
	result := decimal.NewFromInt(1)
	for i := int64(0); i < exponent; i++ {
		result = result.Mul(base)
	}
	return result
}
func FutureValuesOfASeriesFormulaOld(taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64, aporte_primeiro_dia bool) float64 {
    //PMT × {[(1 + r/n)^(nt) - 1] / (r/n)} x (1 + r/n)
    fmt.Println("Args >> ", taxa_juros_decimal, dias_liquidos, anos, valor_aporte, aporte_primeiro_dia)
    fmt.Println("DEBUG", taxa_juros_decimal/dias_liquidos)
	fator_de_crescimento := math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos) - 1
    fmt.Println("Fator de crescimento >> ", fator_de_crescimento)
	fator_de_multiplicacao := fator_de_crescimento / (taxa_juros_decimal / 12)
    fmt.Println("Fator de multiplicaçao >> ", fator_de_multiplicacao)
	valor_futuro := valor_aporte * fator_de_multiplicacao
    if aporte_primeiro_dia {
        return valor_futuro * (1 + (taxa_juros_decimal/12 ))// trocar para a quantidadee de meses restantes para o ano da data inicial acabar))
    }
    return valor_futuro
}
func FutureValuesOfASeriesFormula(taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64, aporte_primeiro_dia bool) (float64, error) {
    taxa_juros_decimal_decimal := decimal.NewFromFloat(taxa_juros_decimal)
    dias_liquidos_decimal := decimal.NewFromFloat(dias_liquidos)
    anos_decimal := decimal.NewFromFloat(anos)
    valor_aporte_decimal := decimal.NewFromFloat(valor_aporte)

    if dias_liquidos_decimal.IsZero() {
        return 0, fmt.Errorf("dias_liquidos cannot be zero")
    }

    fator_crescimento := DecimalPow(
        taxa_juros_decimal_decimal.DivRound(dias_liquidos_decimal, decimalPrecision).Add(decimal.NewFromInt(1)),
        dias_liquidos_decimal.Mul(anos_decimal).IntPart(),
    ).Sub(decimal.NewFromInt(1))

    fator_mutiplicacao := fator_crescimento.DivRound(
        taxa_juros_decimal_decimal.DivRound(decimal.NewFromInt(12), decimalPrecision),
        decimalPrecision,
    )

    valor_futuro := valor_aporte_decimal.Mul(fator_mutiplicacao)
    if aporte_primeiro_dia {
        result := valor_futuro.Mul(
            taxa_juros_decimal_decimal.DivRound(decimal.NewFromInt(12), decimalPrecision).Add(decimal.NewFromInt(1)),
        )
        float_result, _ := result.Round(decimalPrecision).Float64()
        return float_result, nil
    }
    float_result, _ := valor_futuro.Round(decimalPrecision).Float64()
    return float_result, nil
}
func CompoundInterestFormula(valor_inicial, taxa_juros_decimal, dias_liquidos, anos float64) float64 {
//    fmt.Println("Args CIV >> ", valor_inicial, taxa_juros_decimal, dias_liquidos, anos)
    valor_final := valor_inicial * math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos)
    return valor_final
}
func CifAndFvs(valor_inicial, taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64) float64 {
    aux, err := FutureValuesOfASeriesFormula(taxa_juros_decimal, dias_liquidos, anos, valor_aporte, true)
    if err != nil {
        fmt.Println("CifAndFvs")
        panic(err)
    }
    retorno := CompoundInterestFormula(valor_inicial, taxa_juros_decimal, dias_liquidos, anos) + aux
    return retorno
}
type FVSMonthlyMap struct {
    Juros float64 `json:"juros"`
    Acumulado float64 `json:"acumulado"`
    Mes int `json:"mes"`
    JurosFormatado string `json:"juros_formatado"`
    AcumuladoFormatado string `json:"acumulado_formatado"`
    Data time.Time `json:"data"`
    DataMesAno string `json:"data_mes_ano"`
}
func FutureValueOfASeriesMonthly(valor_inicial, taxa_juros_decimal, dias_liquidos, valor_aporte float64, quantidade_meses float64, aporte_primeiro_dia bool, data_inicial time.Time) ([]FVSMonthlyMap, error) {
    valor_acumulado := decimal.NewFromFloat(valor_inicial).Add(decimal.NewFromInt(0))
    var mapa_meses []FVSMonthlyMap
    quantidade_meses_int := int(quantidade_meses)
    //resto_quantidade_meses := quantidade_meses - float64(quantidade_meses_int)
    taxa_mensal :=  decimal.NewFromFloat(taxa_juros_decimal).DivRound(decimal.NewFromInt(12), decimalPrecision)
    aux_date := data_inicial
    valor_aporte_decimal := decimal.NewFromFloat(valor_aporte)
    for mes := 0; mes < quantidade_meses_int; mes++ {
        if aporte_primeiro_dia {
            valor_acumulado = valor_acumulado.Add(valor_aporte_decimal)
        }
        juros := valor_acumulado.Mul(taxa_mensal)
        valor_acumulado = valor_acumulado.Add(juros)
        if !aporte_primeiro_dia {
            valor_acumulado = valor_acumulado.Add(valor_aporte_decimal)
        }
        juros_float, _ := juros.Round(decimalPrecision).Float64()
        acumulado_float, _ := valor_acumulado.Round(decimalPrecision).Float64()
        curr := FVSMonthlyMap{
            Juros: juros_float,
            Acumulado: acumulado_float,
            Mes: mes +1,
            JurosFormatado: core.FormatarValorMonetario(juros_float),
            AcumuladoFormatado: core.FormatarValorMonetario(acumulado_float),
            Data: aux_date,
        }
        mapa_meses = append(mapa_meses, curr)
        aux_date = aux_date.AddDate(0, 1, 0)
    }
    /*valor_acumulado += valor_aporte
    valor_acumulado *= (1+(taxa_mensal*resto_quantidade_meses))*/
    return mapa_meses, nil
}

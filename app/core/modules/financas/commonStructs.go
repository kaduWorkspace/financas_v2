package financas

import (
	"fmt"
	"goravel/app/core"
	"math"
	"time"

	"github.com/shopspring/decimal"
)

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
    //log.Printf("Input values - taxa_juros_decimal: %f, dias_liquidos: %f, anos: %f, valor_aporte: %f, aporte_primeiro_dia: %t",taxa_juros_decimal, dias_liquidos, anos, valor_aporte, aporte_primeiro_dia)
    //log.Printf("Converted taxa_juros_decimal to decimal: %s", taxa_juros_decimal_decimal.String())
    //log.Printf("Converted dias_liquidos to decimal: %s", dias_liquidos_decimal.String())
    //log.Printf("Converted anos to decimal: %s", anos_decimal.String())
    //log.Printf("Converted valor_aporte to decimal: %s", valor_aporte_decimal.String())
    fator_crescimento := DecimalPow(
        taxa_juros_decimal_decimal.Div(dias_liquidos_decimal).Add(decimal.NewFromInt(1)),
        dias_liquidos_decimal.Mul(anos_decimal).IntPart(),
    ).Sub(decimal.NewFromInt(1))
    fator_mutiplicacao := fator_crescimento.Div( taxa_juros_decimal_decimal.Div(decimal.NewFromInt(12)) )

    //fmt.Println("Fator de crescimento >> ", fator_crescimento)
    //fmt.Println("Fator de multiplicaçao >> ", fator_mutiplicacao)
    valor_futuro := valor_aporte_decimal.Mul(fator_mutiplicacao)
    if aporte_primeiro_dia {
        result := valor_futuro.Mul(
            taxa_juros_decimal_decimal.Div(decimal.NewFromInt(12)).Add(decimal.NewFromInt(1)),
        )
        return result.InexactFloat64(), nil
    }
    return valor_futuro.InexactFloat64(), nil
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
    taxa_mensal :=  decimal.NewFromFloat(taxa_juros_decimal).Div(decimal.NewFromInt(12))
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
        curr := FVSMonthlyMap{
            Juros: juros.InexactFloat64(),
            Acumulado: valor_acumulado.InexactFloat64(),
            Mes: mes +1,
            JurosFormatado: core.FormatarValorMonetario(juros.InexactFloat64()),
            AcumuladoFormatado: core.FormatarValorMonetario(valor_acumulado.InexactFloat64()),
            Data: aux_date,
        }
        mapa_meses = append(mapa_meses, curr)
        aux_date = aux_date.AddDate(0, 1, 0)
    }
    /*valor_acumulado += valor_aporte
    valor_acumulado *= (1+(taxa_mensal*resto_quantidade_meses))*/
    return mapa_meses, nil
}

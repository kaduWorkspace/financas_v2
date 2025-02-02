package financas

import (
	"fmt"
	"math"
)
func FutureValuesOfASeriesFormula(taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64, aporte_primeiro_dia bool) float64 {
    //PMT × {[(1 + r/n)^(nt) - 1] / (r/n)} x (1 + r/n)
//    fmt.Println("Args >> ", taxa_juros_decimal, dias_liquidos, anos, valor_aporte, aporte_primeiro_dia)
	fator_de_crescimento := math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos) - 1
	fator_de_multiplicacao := fator_de_crescimento / (taxa_juros_decimal / 12)
	valor_futuro := valor_aporte * fator_de_multiplicacao
    if aporte_primeiro_dia {
        return valor_futuro * (1 + (taxa_juros_decimal/12 ))// trocar para a quantidadee de meses restantes para o ano da data inicial acabar))
    }
    return valor_futuro
}
func CompoundInterestFormula(valor_inicial, taxa_juros_decimal, dias_liquidos, anos float64) float64 {
//    fmt.Println("Args CIV >> ", valor_inicial, taxa_juros_decimal, dias_liquidos, anos)
    valor_final := valor_inicial * math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos)
    return valor_final
}
func CifAndFvs(valor_inicial, taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64) float64 {
    retorno := CompoundInterestFormula(valor_inicial, taxa_juros_decimal, dias_liquidos, anos) + FutureValuesOfASeriesFormula(taxa_juros_decimal, dias_liquidos, anos, valor_aporte, true)
    return retorno
}
func FutureValueOfASeriesMonthly(taxa_juros_decimal, dias_liquidos, valor_aporte float64, quantidade_meses float64, aporte_primeiro_dia bool, taxa_anos float64) []float64 {
    valor_acumulado := 0.0
    var mapa_meses []float64
    quantidade_meses_int := int(quantidade_meses)
    resto_quantidade_meses := quantidade_meses - float64(quantidade_meses_int)
    taxa_mensal := math.Pow(1 + taxa_juros_decimal/12, taxa_anos) - 1
    for mes := 0; mes < quantidade_meses_int; mes++ {
        if aporte_primeiro_dia {
            valor_acumulado += valor_aporte
        }
        valor_acumulado *=  (1+taxa_mensal)
        mapa_meses = append(mapa_meses, valor_acumulado)
        if !aporte_primeiro_dia {
            valor_acumulado += valor_aporte
       }

    }
    fmt.Println("Debug", valor_acumulado, resto_quantidade_meses)
    valor_acumulado += valor_aporte
    valor_acumulado *= (1+(taxa_mensal*resto_quantidade_meses))
    mapa_meses = append(mapa_meses, valor_acumulado)
    fmt.Println("Debug", valor_acumulado)
    return mapa_meses
}

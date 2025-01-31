package financas

import (
	"math"
)
func FutureValuesOfASeriesFormula(taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64, aporte_primeiro_dia bool) float64 {
    //PMT × {[(1 + r/n)^(nt) - 1] / (r/n)} x (1 + r/n)
	fator_de_crescimento := math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos) - 1
	fator_de_multiplicacao := fator_de_crescimento / (taxa_juros_decimal / 12)
	valor_futuro := valor_aporte * fator_de_multiplicacao
    if aporte_primeiro_dia {
        return valor_futuro * (1 + (taxa_juros_decimal/12 ))// trocar para a quantidadee de meses restantes para o ano da data inicial acabar))
    }
    return valor_futuro
}
func CompoundInterestFormula(valor_inicial, taxa_juros_decimal, dias_liquidos, anos float64) float64 {
    valor_final := valor_inicial * math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos)
    return valor_final
}
func CifAndFvs(valor_inicial, taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64) float64 {
    return CompoundInterestFormula(valor_inicial, taxa_juros_decimal, dias_liquidos, anos) + FutureValuesOfASeriesFormula(taxa_juros_decimal, dias_liquidos, anos, valor_aporte, false)
}

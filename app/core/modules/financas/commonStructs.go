package financas

import (
	"fmt"
	"goravel/app/core"
	"math"
	"time"
)
func FutureValuesOfASeriesFormula(taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64, aporte_primeiro_dia bool) float64 {
    //PMT × {[(1 + r/n)^(nt) - 1] / (r/n)} x (1 + r/n)
    fmt.Println("Args >> ", taxa_juros_decimal, dias_liquidos, anos, valor_aporte, aporte_primeiro_dia)
	fator_de_crescimento := math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos) - 1
    fmt.Println("Fator de crescimento >> ", fator_de_crescimento)
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
type FVSMonthlyMap struct {
    Juros float64 `json:"juros"`
    Acumulado float64 `json:"acumulado"`
    Mes int `json:"mes"`
    JurosFormatado string `json:"juros_formatado"`
    AcumuladoFormatado string `json:"acumulado_formatado"`
    Data time.Time `json:"data"`
    DataMesAno string `json:"data_mes_ano"`
}
func FutureValueOfASeriesMonthly(valor_inicial, taxa_juros_decimal, dias_liquidos, valor_aporte float64, quantidade_meses float64, aporte_primeiro_dia bool, data_inicial time.Time) []FVSMonthlyMap {
    valor_acumulado := 0.0 + valor_inicial
    var mapa_meses []FVSMonthlyMap
    quantidade_meses_int := int(quantidade_meses)
    //resto_quantidade_meses := quantidade_meses - float64(quantidade_meses_int)
    taxa_mensal := taxa_juros_decimal/12
    aux_date := data_inicial
    for mes := 0; mes < quantidade_meses_int; mes++ {
        if aporte_primeiro_dia {
            valor_acumulado += valor_aporte
        }
        juros := valor_acumulado * taxa_mensal
        valor_acumulado += juros
        if !aporte_primeiro_dia {
            valor_acumulado += valor_aporte
        }
        curr := FVSMonthlyMap{
            Juros: juros,
            Acumulado: valor_acumulado,
            Mes: mes +1,
            JurosFormatado: core.FormatarValorMonetario(juros),
            AcumuladoFormatado: core.FormatarValorMonetario(valor_acumulado),
            Data: aux_date,
        }
        mapa_meses = append(mapa_meses, curr)
        aux_date = aux_date.AddDate(0, 1, 0)
    }
    /*valor_acumulado += valor_aporte
    valor_acumulado *= (1+(taxa_mensal*resto_quantidade_meses))*/
    return mapa_meses
}

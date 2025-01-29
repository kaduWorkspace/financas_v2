package financas

import (
	"errors"
	"fmt"
	"goravel/app/core"
	"math"
	"time"

	"github.com/goravel/framework/facades"
)

// Sempre será pré fixado, pois a taxa é fixada antes do inicio do investimento e se mantem ate o final.
type SimulacaoCDB struct {
    dataInicial time.Time `json:"data_inicial"`
    dataFinal time.Time `json:"data_final"`
    ValorEstimadoFuturo float64 `json:"valor_estimado_futuro"`
    valorInicial float64
    fracaoAnos float64
    valorAporte float64
    taxaJurosDecimal float64
    diasLiquidez float64
}
func (self *SimulacaoCDB) SetDatas(data_inicial string, data_final string) error {
    data_layout := "2006-01-02"
    data_inicial_obj, err := time.Parse(data_layout, data_inicial)
    if err != nil {
        facades.Log().Debug(err.Error())
        return err
    }
    data_final_obj, err := time.Parse(data_layout, data_final)
    if err != nil {
        facades.Log().Debug(err.Error())
        return err
    }
    self.dataInicial = data_inicial_obj
    self.dataFinal = data_final_obj
    return nil
}
func (self *SimulacaoCDB) SetValorAporte(aporte float64) {
    self.valorAporte = aporte
}
func (self *SimulacaoCDB) SetDiasLiquidezPorAno (quantidade float64) {
    self.diasLiquidez = quantidade
}
func (self *SimulacaoCDB) JurosComposto() {}
func (self *SimulacaoCDB) SetDiasDeLiquidesPorAno(quantidade int) {
    self.diasLiquidez = float64(quantidade)
}
func (self *SimulacaoCDB) SetTaxaDeJurosDecimal(valor float64, tipo string) error {
    if tipo == "porcento anual" {
        self.taxaJurosDecimal = valor / 100
        return nil
    }
    return errors.New("Tipo não é valido!")
}
func (self *SimulacaoCDB) SetTaxaAnosApartirPeriodoDeDatas() error {
    mapa_dias_por_ano := map[int]int {}
    mapa_dias_por_ano[self.dataInicial.Year()] = int(math.Abs(float64(self.dataInicial.YearDay() - core.DiasNoAnoV2(self.dataInicial))))
    aux_date := self.dataInicial
    aux_date = aux_date.AddDate(1,0,0)
    for aux_date.Year() <= self.dataFinal.Year() {
        if self.dataFinal.Year() == aux_date.Year() {
            mapa_dias_por_ano[self.dataFinal.Year()] = aux_date.YearDay()
            break
        } else {
            mapa_dias_por_ano[aux_date.Year()] = core.DiasNoAnoV2(aux_date)
        }
        aux_date = aux_date.AddDate(1,0,0)
        if aux_date.Year() != self.dataFinal.Year() {
            aux_date = time.Date(aux_date.Year(), time.January, 0, 0, 0, 0, 0, time.UTC)
        } else {
            aux_date = self.dataFinal
        }
    }
    quantidade_anos := 0.0
    dateLayout := "2006-01-02"
    for ano, dias := range mapa_dias_por_ano {
        aux_date_2, err := time.Parse(dateLayout, fmt.Sprintf("%d-01-01",ano))
        if err != nil {
            return err
        }
        total_dias_ano := core.DiasRestantesNoAno(aux_date_2)
        quantidade_anos = quantidade_anos + float64(dias)/float64(total_dias_ano)
    }
    self.fracaoAnos = quantidade_anos
    return nil
}
func (self *SimulacaoCDB) SetValorInicial(valor_inicial float64) {
    self.valorInicial = valor_inicial
}
func (self *SimulacaoCDB) JurosCompostoComAporteMensal() float64 {
    //PMT × {[(1 + r/n)^(nt) - 1] / (r/n)} x (1 + r/n)
	fator_de_crescimento := math.Pow(1 + (self.taxaJurosDecimal/self.diasLiquidez), self.diasLiquidez*self.fracaoAnos) - 1
	fator_de_multiplicacao := fator_de_crescimento / (self.taxaJurosDecimal / self.diasLiquidez)
	valor_futuro := self.valorAporte * fator_de_multiplicacao
    return valor_futuro * (1 + (self.taxaJurosDecimal/12))
}
func (self *SimulacaoCDB) GetDiasDeLiquidezPorAno() float64 {
    return self.diasLiquidez
}
func (self *SimulacaoCDB) GetTaxaDeJurosDecimal() float64 {
    return self.taxaJurosDecimal
}
func (self *SimulacaoCDB) GetTaxaAnos() float64 {
    return self.fracaoAnos
}
func (self *SimulacaoCDB) GetDataInicial() time.Time {
    return self.dataInicial
}
func (self *SimulacaoCDB) GetDataFinal() time.Time {
    return self.dataFinal
}
func (self *SimulacaoCDB) GetValorAporte() float64 {
    return self.valorAporte
}
func (self *SimulacaoCDB) GetTaxaSelic() float64 {
    valor_selic := 11.25 //padrao
    result, err := core.HttpRequest("https://www.bcb.gov.br/api/servico/sitebcb//taxaselic/ultima?withCredentials=true", "GET", map[string]string{"content-type":"text/plain"}, "")
    if err == nil {
        type RetornoBancoCentralApi struct {
            MetaSelic          float64 `json:"MetaSelic"`
            DataReuniaoCopom   string  `json:"DataReuniaoCopom"`
            Vies               string  `json:"Vies"`
        }
        type RetornoBancoCentralApiWrapper struct {
            Conteudo []RetornoBancoCentralApi `json:"conteudo"`
        }
        var bodyRes RetornoBancoCentralApiWrapper
        if err := core.ConverterJson(result, &bodyRes); err == nil && len(bodyRes.Conteudo) != 0 {
            valor_selic = bodyRes.Conteudo[0].MetaSelic
        }
    }
    return valor_selic
}
func (self *SimulacaoCDB) GetValorInicial() float64 {
    return self.valorInicial
}
func FutureValuesOfASeriesFormula(taxa_juros_decimal, dias_liquidos, anos, valor_aporte float64) float64 {
    //PMT × {[(1 + r/n)^(nt) - 1] / (r/n)} x (1 + r/n)
	fator_de_crescimento := math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos) - 1
	fator_de_multiplicacao := fator_de_crescimento / (taxa_juros_decimal / 12)
	valor_futuro := valor_aporte * fator_de_multiplicacao
    return valor_futuro * (1 + (taxa_juros_decimal/12 ))// trocar para a quantidadee de meses restantes para o ano da data inicial acabar))
}
func CompoundInterestFormula(valor_inicial, taxa_juros_decimal, dias_liquidos, anos float64) float64 {
    valor_final := valor_inicial * math.Pow(1 + (taxa_juros_decimal/dias_liquidos), dias_liquidos*anos)
    return valor_final
}

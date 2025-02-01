package financas

import (
	"errors"
	"fmt"
	"goravel/app/core"
	"math"
	"time"

	"github.com/goravel/framework/facades"
)
type TIPO_TAXA_JUROS int
const (
    PROCENTO_ANUAL TIPO_TAXA_JUROS = iota
)
// Sempre será pré fixado, pois a taxa é fixada antes do inicio do investimento e se mantem ate o final.
type SimularJurosComposto struct {
    dataInicial time.Time `json:"data_inicial"`
    dataFinal time.Time `json:"data_final"`
    valorInicial float64
    fracaoAnos float64
    valorAporte float64
    taxaJurosDecimal float64
    diasLiquidez float64
    valorGasto float64
    valorJuros float64
}
func (self *SimularJurosComposto) SetDatas(data_inicial string, data_final string) error {
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
func (self *SimularJurosComposto) SetValorAporte(aporte float64) {
    self.valorAporte = aporte
}
func (self *SimularJurosComposto) SetDiasLiquidezPorAno (quantidade float64) {
    self.diasLiquidez = quantidade
}
func (self *SimularJurosComposto) SetDiasDeLiquidesPorAno(quantidade int) {
    self.diasLiquidez = float64(quantidade)
}
func (self *SimularJurosComposto) SetTaxaDeJurosDecimal(valor float64, tipo TIPO_TAXA_JUROS) error {
    if tipo == PROCENTO_ANUAL {
        self.taxaJurosDecimal = valor / 100
        return nil
    }
    return errors.New("Tipo não é valido!")
}
func (self *SimularJurosComposto) SetTaxaAnosApartirPeriodoDeDatas() error {
    if self.dataFinal.Year() == self.dataInicial.Year() {
        total_dias_ano := core.DiasNoAnoV2(self.dataFinal)
        self.fracaoAnos = (float64(self.dataFinal.YearDay()) - float64(self.dataInicial.YearDay())) / float64(total_dias_ano)
        //fmt.Println("Setando fracao anos para", self.fracaoAnos)
        return nil
    }
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
            aux_date = aux_date.AddDate(1,0,0)
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
func (self *SimularJurosComposto) SetValorInicial(valor_inicial float64) {
    self.valorInicial = valor_inicial
}
func (self *SimularJurosComposto) GetDiasDeLiquidezPorAno() float64 {
    return self.diasLiquidez
}
func (self *SimularJurosComposto) GetTaxaDeJurosDecimal() float64 {
    return self.taxaJurosDecimal
}
func (self *SimularJurosComposto) GetTaxaAnos() float64 {
    return self.fracaoAnos
}
func (self *SimularJurosComposto) GetDataInicial() time.Time {
    return self.dataInicial
}
func (self *SimularJurosComposto) GetDataFinal() time.Time {
    return self.dataFinal
}
func (self *SimularJurosComposto) GetValorAporte() float64 {
    return self.valorAporte
}
func (self *SimularJurosComposto) GetTaxaSelic() float64 {
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
func (self *SimularJurosComposto) GetValorInicial() float64 {
    return self.valorInicial
}
func (self *SimularJurosComposto) SetValorJurosRendido(valor_final float64) {
    if self.GetValorGasto() == 0.0 {
        self.SetValorGasto()
    }
    valor_juros := valor_final - self.GetValorGasto()
    self.valorJuros = valor_juros
}
func (self *SimularJurosComposto) GetValorJurosRendido() float64 {
    return self.valorJuros
}
func (self *SimularJurosComposto) SetValorGasto() {
    valor_gasto := 0.0
    if self.GetValorInicial() > 0 {
        valor_gasto += self.GetValorInicial()
    }
    if self.GetValorAporte() > 0 {
        quantidade_meses := core.MesesEntreDatas(self.GetDataFinal(), self.GetDataInicial())
        valor_gasto = valor_gasto + (self.GetValorAporte() * float64(quantidade_meses))
    }
    self.valorGasto = valor_gasto
}
func (self *SimularJurosComposto) GetValorGasto() float64 {
    return self.valorGasto
}

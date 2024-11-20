package financas

import (
	"goravel/app/core"
	"goravel/app/http/requests"
	"time"
)
type TrackerDiario struct {
    Valorizacao float64 `json:"valorizacao"`
    Data string `json:"data"`
    ResultadoComValorizacao float64 `json:"resultado_com_valorizacao"`
}

func newTrackerDiario(dataInicial time.Time, valorizacao float64, resultadoComValorizacao float64) TrackerDiario {
    return TrackerDiario {
        Data: core.FormataData(dataInicial),
        ResultadoComValorizacao: resultadoComValorizacao,
        Valorizacao: valorizacao,
    }
}
type TrackerMensal struct {
    Dias []TrackerDiario `json:"dias"`
    Valorizacao float64 `json:"valorizacao"`
    DataInicial string `json:"data_inicial"`
    DataFinal string `json:"data_final"`
    ResultadoComValorizacao float64 `json:"resultado_com_valorizacao"`
    Aporte float64 `json:"valor_aporte"`
}
func newTrackerMensal(dataInicial time.Time, aporte float64) TrackerMensal {
    return TrackerMensal {
        DataInicial: core.FormataData(dataInicial),
        DataFinal: "",
        Aporte: aporte,
    }
}
func (self *TrackerMensal) AdicionaDia(dia TrackerDiario) {
    self.Dias = append(self.Dias, dia)
}
func (self *TrackerMensal) finalizarMes(dataFinal time.Time, resultado float64) {
    self.ResultadoComValorizacao = resultado
    self.DataFinal = core.FormataData(dataFinal)
}
type TrackerSemestral struct {
    Meses []TrackerMensal `json:"meses"`
    Valorizacao float64 `json:"valorizacao"`
    DataInicial string `json:"data_inicial"`
    DataFinal string `json:"data_final"`
    ResultadoComValorizacao float64 `json:"resultado_com_valorizacao"`
    Aporte float64 `json:"valor_aporte"`
    Gasto float64 `json:"gasto"`
}
func newTrackerSemestral(dataInicial time.Time, valorAporteSemestral float64) TrackerSemestral {
    return TrackerSemestral {
        DataInicial: core.FormataData(dataInicial),
        Gasto: 0,
        Valorizacao: 0,
        Aporte: valorAporteSemestral,
    }
}
func (self *TrackerSemestral) adicionaGasto(gasto float64) {
    self.Gasto += gasto
}
func (self *TrackerSemestral) adicionaMes(mes TrackerMensal) {
    self.Meses = append(self.Meses, mes)
}
func (self *TrackerSemestral) finalizarSemestre(dataFinal time.Time , resultado float64) {
    self.ResultadoComValorizacao = resultado
    self.DataFinal = core.FormataData(dataFinal)
}
type TrackerAnual struct {
    Semestres []TrackerSemestral `json:"semestres"`
    Valorizacao float64 `json:"valorizacao"`
    DataInicial string `json:"data_inicial"`
    DataFinal string `json:"data_final"`
    ResultadoComValorizacao float64 `json:"resultado_com_valorizacao"`
    Gasto float64 `json:"gasto"`
}
func newTrackerAnual(dataInicial time.Time) TrackerAnual {
    return TrackerAnual {
        DataInicial: core.FormataData(dataInicial),
        Gasto: 0,
        Valorizacao: 0,
    }
}
func (self *TrackerAnual) adicionaGasto(gasto float64) {
    self.Gasto += gasto
}
func (self * TrackerAnual) adicionaSemestre(tracker TrackerSemestral) {
    self.Semestres = append(self.Semestres, tracker)
}
func (self * TrackerAnual) finalizarAno(dataFinal time.Time, resultado float64) {
    self.ResultadoComValorizacao = resultado
    self.DataFinal = core.FormataData(dataFinal)
}
type CalcularJuros struct {
    ValorInicial          float64 `json:"valor_inicial" form:"valor_inicial"`
    DataInicial           string `json:"data_inicial" form:"data_inicial"`
    DataFinal             string `json:"data_final" form:"data_final"`
    AporteMensal          float64  `json:"aporte_mensal" form:"aporte_mensal"`
    AporteSemestral       float64  `json:"aporte_semestral" form:"aporte_semestral"`
    TipoFrequenciaAumentoAporte  string `json:"tipo_frequencia_aumento_aporte" form:"tipo_frequencia_aumento_aporte"`
    ValorAumentoAporte    float64  `json:"valor_aumento_aporte" form:"valor_aumento_aporte"`
    TaxaSelic float64 `json:"taxa_selic_diario"`
    dataInicial time.Time
    dataFinal time.Time
    TrackerAnual TrackerAnual `json:"tracker_anual"`
}
func New(input requests.PostCalcularJuros) (CalcularJuros, error) {
    var self CalcularJuros
    self.DataInicial = input.DataInicial
    self.DataFinal = input.DataFinal
    self.AporteMensal = input.AporteMensal
    self.AporteSemestral = input.AporteSemestral
    self.TipoFrequenciaAumentoAporte = input.TipoFrequenciaAumentoAporte
    self.ValorAumentoAporte = input.ValorAumentoAporte
    self.ValorInicial = input.ValorInicial
    timeInicial, err := time.Parse("02/01/2006", self.DataInicial)
    if err != nil {
        return self, nil
    }
    self.dataInicial = timeInicial
    timeFinal, err := time.Parse("02/01/2006", self.DataFinal)
    if err != nil {
        return self, nil
    }
    self.dataFinal = timeFinal
    self.setTaxaSelic()
    return self, nil
}
func (self *CalcularJuros) setTaxaSelic() error {
    result, err := core.HttpRequest("https://www.bcb.gov.br/api/servico/sitebcb//taxaselic/ultima?withCredentials=true", "GET", map[string]string{"content-type":"text/plain"}, "")
    if err != nil {
        taxaAnual := 11.25
        self.TaxaSelic = float64((taxaAnual / 365) / 100)
        return err
    }
    type RetornoBancoCentralApi struct {
        MetaSelic          float64 `json:"MetaSelic"`
        DataReuniaoCopom   string  `json:"DataReuniaoCopom"`
        Vies               string  `json:"Vies"`
    }
    type RetornoBancoCentralApiWrapper struct {
        Conteudo []RetornoBancoCentralApi `json:"conteudo"`
    }
    var bodyRes RetornoBancoCentralApiWrapper
    if err := core.ConverterJson(result, &bodyRes); err != nil {
        taxaAnual := 11.25
        self.TaxaSelic = float64((taxaAnual / 365) / 100)
    }
    if len(bodyRes.Conteudo) == 0 {
        taxaAnual := 11.25
        self.TaxaSelic = float64((taxaAnual / 365) / 100)
    } else {
        taxaAnual := bodyRes.Conteudo[0].MetaSelic
        self.TaxaSelic = taxaAnual / 365 / 100
    }
    return nil
}
type ResultadoSimulacao struct {
    Anos []TrackerAnual `json:"anos"`
    Valorizacao float64 `json:"valorizacao"`
    ValorFinal float64 `json:"valor_final"`
    ValorInicial float64 `json:"valor_inicial"`
    Gasto float64 `json:"gastos"`
    Diferenca float64 `json:"diferenca"`
    DataInicial string `json:"data_final"`
    DataFinal string `json:"data_inicial"`
}
func newResultadoSimulacao(dataInicial time.Time, dataFinal time.Time, valorInicial float64) ResultadoSimulacao {
    return ResultadoSimulacao {
        DataInicial: core.FormataData(dataInicial),
        DataFinal: core.FormataData(dataFinal),
        ValorInicial: valorInicial,
    }
}
func (self *ResultadoSimulacao) adicionaGasto(gasto float64) {
    self.Gasto += gasto
}

func (self *ResultadoSimulacao) adicionaValorizacao(valorizacao float64) {
    self.Valorizacao += valorizacao
}
func (self *ResultadoSimulacao) finalizarSimulacao(dataFinal time.Time, resultado float64, anos []TrackerAnual) {
    self.Anos = anos
    self.DataFinal = core.FormataData(dataFinal)
    self.ValorFinal = resultado
    self.calcularDiferenca()
}
func (self *ResultadoSimulacao) calcularDiferenca() {
    self.Diferenca = self.ValorFinal - self.ValorInicial
}
func (self *CalcularJuros) Calcular() ResultadoSimulacao {
    resultadoSimulacao := newResultadoSimulacao(self.dataInicial, self.dataFinal, self.ValorInicial)

    dataAtual := self.dataInicial

    resultado := self.ValorInicial + self.AporteMensal

    prxMes := dataAtual.AddDate(0, 1, 0)
    prxSemestre := dataAtual.AddDate(0, 6, 0)
    prxAno := dataAtual.AddDate(1, 0, 0)

    var trackerAnuais []TrackerAnual

    trackerMensalAtual := newTrackerMensal(dataAtual, self.AporteMensal)
    trackerSemestralAtual := newTrackerSemestral(dataAtual, self.AporteSemestral)
    trackerAnualAtual :=  newTrackerAnual(dataAtual)

    trackerAnualAtual.adicionaGasto(self.AporteMensal)
    trackerSemestralAtual.adicionaGasto(self.AporteMensal)
    resultadoSimulacao.adicionaGasto(self.AporteMensal)
    for !dataAtual.After(self.dataFinal) {
        dataAtual = dataAtual.AddDate(0, 0, 1)
        valorizacao := self.TaxaSelic * resultado
        resultadoSimulacao.adicionaValorizacao(valorizacao)
        resultado = resultado + valorizacao
        trackerMensalAtual.Valorizacao += valorizacao
        trackerSemestralAtual.Valorizacao += valorizacao
        trackerAnualAtual.Valorizacao += valorizacao
        trackerDiario := newTrackerDiario(dataAtual, valorizacao, resultado)
        trackerMensalAtual.AdicionaDia(trackerDiario)
        if dataAtual.Equal(prxMes) {
            resultado += self.AporteMensal
            trackerMensalAtual.finalizarMes(dataAtual, resultado)

            trackerSemestralAtual.adicionaGasto(self.AporteMensal)
            trackerSemestralAtual.adicionaMes(trackerMensalAtual)

            trackerAnualAtual.adicionaGasto(self.AporteMensal)
            trackerMensalAtual = newTrackerMensal(dataAtual, self.AporteMensal)
            prxMes = prxMes.AddDate(0, 1, 0)

            resultadoSimulacao.adicionaGasto(self.AporteMensal)
        }
        if dataAtual.Equal(prxSemestre) {
            if self.TipoFrequenciaAumentoAporte == "semestral" {
                self.AporteMensal += self.ValorAumentoAporte
            }
            resultado += self.AporteSemestral
            trackerSemestralAtual.adicionaGasto(self.AporteSemestral)
            trackerSemestralAtual.finalizarSemestre(dataAtual, resultado)

            trackerAnualAtual.adicionaGasto(self.AporteSemestral)
            trackerAnualAtual.adicionaSemestre(trackerSemestralAtual)

            trackerSemestralAtual = newTrackerSemestral(dataAtual, self.AporteSemestral)
            prxSemestre = prxSemestre.AddDate(0, 6, 0)

            resultadoSimulacao.adicionaGasto(self.AporteSemestral)
        }
        if dataAtual.Equal(prxAno) {
            if self.TipoFrequenciaAumentoAporte == "anual" {
                self.AporteMensal += self.ValorAumentoAporte
            }
            trackerAnualAtual.finalizarAno(dataAtual, resultado)
            trackerAnuais = append(trackerAnuais, trackerAnualAtual)

            trackerAnualAtual = newTrackerAnual(dataAtual)
            prxAno = prxAno.AddDate(1, 0, 0)
        }
        if dataAtual.Equal(self.dataFinal) {
            trackerMensalAtual.finalizarMes(dataAtual, resultado)

            trackerSemestralAtual.adicionaMes(trackerMensalAtual)
            trackerSemestralAtual.finalizarSemestre(dataAtual, resultado)

            trackerAnualAtual.adicionaSemestre(trackerSemestralAtual)
            trackerAnualAtual.finalizarAno(dataAtual, resultado)

            trackerAnuais = append(trackerAnuais, trackerAnualAtual)
        }
    }
    resultadoSimulacao.finalizarSimulacao(dataAtual, resultado, trackerAnuais)
    return resultadoSimulacao
}

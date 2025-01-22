package financas

import (
	"bytes"
	"encoding/json"
	"goravel/app/core"
	"goravel/app/http/requests"
	"image/color"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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
    taxaAnual float64
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
        self.taxaAnual = 11.25
        self.TaxaSelic = float64((self.taxaAnual / 365) / 100)
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
        self.taxaAnual = 11.25
        self.TaxaSelic = float64((self.taxaAnual / 365) / 100)
    }
    if len(bodyRes.Conteudo) == 0 {
        self.taxaAnual = 11.25
        self.TaxaSelic = float64((self.taxaAnual / 365) / 100)
    } else {
        self.taxaAnual = bodyRes.Conteudo[0].MetaSelic
        self.TaxaSelic = self.taxaAnual / 365 / 100
    }
    return nil
}
type ResultadoSimulacao struct {
    Anos []TrackerAnual `json:"anos"`
    Meses []TrackerMensal `json:"meses"`
    Semestres []TrackerSemestral `json:"semestres"`
    Dias []TrackerDiario `json:"dias"`
    Valorizacao float64 `json:"valorizacao"`
    ValorFinal float64 `json:"valor_final"`
    ValorInicial float64 `json:"valor_inicial"`
    Gasto float64 `json:"gastos"`
    Diferenca float64 `json:"diferenca"`
    DataInicial string `json:"data_inicial"`
    DataFinal string `json:"data_final"`
}
func (self *ResultadoSimulacao) ToJson() (string, error) {
    jsonData, err := json.Marshal(self)
    if err != nil {
        return "", err
    }
    return string(jsonData), nil
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
func (self *ResultadoSimulacao) SetMeses() {
    for _, ano := range self.Anos {
        for _, semestre := range ano.Semestres {
            for _, mes := range semestre.Meses {
                self.Meses = append(self.Meses, TrackerMensal{
                    Valorizacao: mes.Valorizacao,
                    ResultadoComValorizacao: mes.ResultadoComValorizacao,
                    Aporte: mes.Aporte,
                    DataInicial: mes.DataInicial,
                    DataFinal: mes.DataFinal,
                })
            }
        }
    }
}
func (self *ResultadoSimulacao) SetDias() {
    for _, ano := range self.Anos {
        for _, semestre := range ano.Semestres {
            for _, mes := range semestre.Meses {
                for _, dia := range mes.Dias {
                    self.Dias = append(self.Dias, TrackerDiario{
                        Valorizacao: dia.Valorizacao,
                        ResultadoComValorizacao: dia.ResultadoComValorizacao,
                        Data: dia.Data,
                    })
                }
            }
        }
    }
}
func (self *ResultadoSimulacao) Plot(tipo string) error {
    if tipo == "mes" {
        var valorizacoes []float64
        var datas []string
        for _, mes := range self.Meses {
            valorizacoes = append(valorizacoes, mes.Valorizacao)
            datas = append(datas, mes.DataFinal)
        }
        p := plot.New()
        p.Title.Text = "Valorizacao dos dias"
        p.X.Label.Text = "Dias"
        p.Y.Label.Text = "Valorização"
        p.Add(plotter.NewGrid())

        pontos := make(plotter.XYs, len(valorizacoes))
        for i, valor := range valorizacoes {
            pontos[i].X = float64(i)
            pontos[i].Y = valor
        }
        linha, err := plotter.NewLine(pontos)
        if err != nil {
            return err
        }
        linha.Color = color.RGBA{R: 255, A: 255}
        linha.Width = vg.Points(2)
        p.Add(linha)
        p.NominalX(datas...)
        if err := p.Save(100*vg.Centimeter, 50*vg.Centimeter, "./public/graficos/grafico_meses.png"); err != nil {
            return err
        }
    }
    return nil
}
type Graficos struct {
    MesesValorizacoes string
    DiasValorizacoes string
    AnosValorizacoes string
    SemestresValorizacoes string
}

func (self *ResultadoSimulacao) Chart() (error, Graficos) {
    var resultado Graficos
    graficoMeses := charts.NewBar()
    var dadosGraficoMes []opts.BarData
    var meses []string
    for _, mes := range self.Meses {
        dadosGraficoMes = append(dadosGraficoMes, opts.BarData{
            Value: mes.Valorizacao,
        })
        meses = append(meses, mes.DataFinal)
    }

    graficoMeses.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "",
        }),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Mês",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Valorização",
		}),
    )
    graficoMeses.SetXAxis(meses).AddSeries("Valorização primeiro set", dadosGraficoMes)
    var buffer bytes.Buffer
    if err :=  graficoMeses.Render(&buffer);err != nil {
        println("Não foi possivel criar o gráfico de valorizacao mensais")
        return err, resultado
    }
    resultado.MesesValorizacoes = buffer.String()
    println("Grafico mêses Criado!")

    graficoDias := charts.NewBar()
    var dadosGraficoDia []opts.BarData
    var dias []string
    for _, dia := range self.Dias {
        dadosGraficoDia = append(dadosGraficoDia, opts.BarData{
            Value: dia.Valorizacao,
        })
        dias = append(dias, dia.Data)
    }

    graficoDias.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "",
        }),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Dia",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Valorização",
		}),
    )
    graficoDias.SetXAxis(dias).AddSeries("Valorização do dia", dadosGraficoDia)
    var bufferDias bytes.Buffer
    if err :=  graficoDias.Render(&bufferDias);err != nil {
        println("Não foi possivel criar o gráfico de valorizacao diarias")
        return err, resultado
    }
    resultado.DiasValorizacoes = buffer.String()
    println("Grafico mêses Criado!")


    graficoSemestres := charts.NewBar()
    var dadosGraficoSemestre []opts.BarData
    var semestres []string
    for _, semestre := range self.Semestres {
        dadosGraficoSemestre = append(dadosGraficoSemestre, opts.BarData{
            Value: semestre.Valorizacao,
        })
        semestres = append(semestres, semestre.DataInicial)
    }

    graficoSemestres.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "",
        }),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Semestre",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Valorização",
		}),
    )
    graficoSemestres.SetXAxis(semestres).AddSeries("Valorização dos semestres", dadosGraficoSemestre)
    var bufferSemestres bytes.Buffer
    if err :=  graficoSemestres.Render(&bufferSemestres);err != nil {
        println("Não foi possivel criar o gráfico de valorizacao semestres")
        return err, resultado
    }
    resultado.SemestresValorizacoes = buffer.String()
    println("Grafico mêses Criado!")

    graficoAnos := charts.NewBar()
    var dadosGraficoAno []opts.BarData
    var anos []string
    for _, ano := range self.Anos {
        dadosGraficoAno = append(dadosGraficoAno, opts.BarData{
            Value: ano.Valorizacao,
        })
        anos = append(anos, ano.DataInicial)
    }

    graficoAnos.SetGlobalOptions(
        charts.WithTitleOpts(opts.Title{
            Title: "",
        }),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Ano",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Valorização",
		}),
    )
    graficoAnos.SetXAxis(anos).AddSeries("Valorização dos anos", dadosGraficoAno)
    var bufferAnos bytes.Buffer
    if err :=  graficoAnos.Render(&bufferAnos);err != nil {
        println("Não foi possivel criar o gráfico de valorizacao anos")
        return err, resultado
    }
    resultado.AnosValorizacoes = buffer.String()
    println("Grafico anos Criado!")
    return nil, resultado
}
func (self *ResultadoSimulacao) SetSemestres() {
    for _, ano := range self.Anos {
        for _, semestre := range ano.Semestres {
            self.Semestres = append(self.Semestres, TrackerSemestral{
                Valorizacao: semestre.Valorizacao,
                ResultadoComValorizacao: semestre.ResultadoComValorizacao,
                Aporte: semestre.Aporte,
                DataInicial: semestre.DataInicial,
                DataFinal: semestre.DataFinal,
                Gasto: semestre.Gasto,
            })
        }
    }
}
func (self *ResultadoSimulacao) adicionaValorizacao(valorizacao float64) {
    self.Valorizacao += valorizacao
}
func (self *ResultadoSimulacao) finalizarSimulacao(resultado float64, anos []TrackerAnual) {
    self.Anos = anos
    self.ValorFinal = resultado
    self.calcularDiferenca()
}
func (self *ResultadoSimulacao) calcularDiferenca() {
    self.Diferenca = self.ValorFinal - self.Gasto
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

    /*trackerAnualAtual.adicionaGasto(self.AporteMensal)
    trackerSemestralAtual.adicionaGasto(self.AporteMensal)
    resultadoSimulacao.adicionaGasto(self.AporteMensal)*/
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
    resultadoSimulacao.finalizarSimulacao(resultado, trackerAnuais)
    return resultadoSimulacao
}

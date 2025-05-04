package Feature

import (
	"fmt"
	"testing"

	"goravel/app/core/modules/financas"
	"goravel/tests"

	"github.com/stretchr/testify/suite"
)

type FutureValueOfASeriesSuite struct {
	suite.Suite
	tests.TestCase
}

func TestFutureValueOfASeriesSuite(t *testing.T) {
	suite.Run(t, new(FutureValueOfASeriesSuite))
}

// SetupTest will run before each test in the suite.
func (s *FutureValueOfASeriesSuite) SetupTest() {
}

// TearDownTest will run after each test in the suite.
func (s *FutureValueOfASeriesSuite) TearDownTest() {
}
/*func (self *FutureValueOfASeriesSuite) TestFutureValueOfASeries() {
    servico := financas.SimularJurosComposto {}
    servico.SetDiasDeLiquidesPorAno(255)
    servico.SetTaxaDeJurosDecimal(13.25, financas.PORCENTO_ANUAL)
    servico.SetValorAporte(833.00)
    servico.SetDiasDeLiquidesPorAno(12)
    valorizacao, _ := financas.FutureValuesOfASeriesFormula(servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), 1, servico.GetValorAporte(), false)
    valor_alvo := 10625.956528507051
    if valorizacao != valor_alvo {
        self.Fail(fmt.Sprintf("[FVS] Valorização deveria ser %f, retornado %f", valor_alvo, valorizacao))
    }
}*/
/*
func (self *FutureValueOfASeriesSuite) TestFVSWithPrecision() {
    servico := financas.SimularJurosComposto {}
    data_inicial := "2025-02-01"
    data_final := "2037-02-01"
    fmt.Println("Datas inicial e final >> ", data_inicial, data_final)
    servico.SetDatas(data_inicial, data_final)
    servico.SetTaxaAnosApartirPeriodoDeDatas()
    servico.SetTaxaAnos(float64(int(servico.GetTaxaAnos())))
    servico.SetDiasDeLiquidesPorAno(12)
    servico.SetTaxaDeJurosDecimal(13.25, financas.PORCENTO_ANUAL)
    servico.SetValorAporte(1833.00)
    resultOld := financas.FutureValuesOfASeriesFormulaOld(servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), servico.GetTaxaAnos(), servico.GetValorAporte(), true)
    resulPrecision, err := financas.FutureValuesOfASeriesFormula(servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), servico.GetTaxaAnos(), servico.GetValorAporte(), true)
    fmt.Println("Resultado old >> ", resultOld)
    fmt.Println("Resultado >> ", resulPrecision)
    if err != nil {
        panic(err)
    }
    self.Fail("Stop")
}*/
func (self *FutureValueOfASeriesSuite) TestMonthlyFVS() {
    servico := financas.SimularJurosComposto {}
    data_inicial := "2025-02-01"
    data_final := "2046-02-01"
    fmt.Println("Datas inicial e final >> ", data_inicial, data_final)
    servico.SetDatas(data_inicial, data_final)
    servico.SetTaxaAnosApartirPeriodoDeDatas()
    servico.SetTaxaAnos(float64(int(servico.GetTaxaAnos())))
    servico.SetDiasDeLiquidesPorAno(12)
    servico.SetTaxaDeJurosDecimal(13.25, financas.PORCENTO_ANUAL)
    servico.SetValorAporte(833.00)
    reusultado_padrao, err := financas.FutureValuesOfASeriesFormula(servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), servico.GetTaxaAnos(), servico.GetValorAporte(), true)
    if err != nil {
        panic(err)
    }
    servico.SetTaxaMeses(servico.GetTaxaAnos() * 12)
    mapas, err := financas.FutureValueOfASeriesMonthly(servico.GetValorInicial(), servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), servico.GetValorAporte(), servico.GetTaxaAnos() * 12, true, servico.GetDataInicial())
    if err != nil {
        self.Fail(err.Error())
    }
    //fmt.Println("Taxa anos", servico.GetTaxaAnos())
    //fmt.Println("Taxa Meses", servico.GetTaxaMeses())
    //fmt.Println("Resultado padrao >> ", reusultado_padrao)
    //fmt.Println("Ultimo valor >> ", mapas[len(mapas) -1])
    ultimo_valor := mapas[len(mapas) -1]
    if int(reusultado_padrao) != int(ultimo_valor.Acumulado) {
        self.Fail(fmt.Sprintf("[FVS] Valorização deveria ser %f, retornado %f", reusultado_padrao, ultimo_valor.Acumulado))
    }
    fmt.Println("Esperado", reusultado_padrao, "Recebido", ultimo_valor.Acumulado)

}
func (self *FutureValueOfASeriesSuite) TestCompoundInterestFomula() {
    servico := financas.SimularJurosComposto {}
    servico.SetDiasDeLiquidesPorAno(255)
    servico.SetTaxaDeJurosDecimal(13.25, financas.PORCENTO_ANUAL)
    servico.SetValorInicial(9500.00)
    servico.SetDiasDeLiquidesPorAno(12)
    valorizacao := financas.CompoundInterestFormula(servico.GetValorInicial(), servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), 1)
    valor_alvo := 10838.077509029437
    if valorizacao != valor_alvo {
        self.Fail(fmt.Sprintf("[CI] Valorização deveria ser %f, retornado %f", valor_alvo, valorizacao))
    }
}
func (self *FutureValueOfASeriesSuite) TestTaxaAnos() {
    servico := financas.SimularJurosComposto {}
    data_inicial := "2025-02-01"
    data_final := "2026-04-01"
    servico.SetDatas(data_inicial, data_final)
    servico.SetTaxaAnosApartirPeriodoDeDatas()
    valor_alvo := 1.1648351648351647
    if servico.GetTaxaAnos() != valor_alvo {
        fmt.Println("Valor alvo nao atingido >> ", valor_alvo, "Valor atingido", servico.GetTaxaAnos(), "Datas", servico.GetDataInicial(), servico.GetDataFinal())
        self.T().Fail()
    }

}
func (self *FutureValueOfASeriesSuite) TestFutureValueOfASeriesMonthly() {
    data_inicial := "2025-02-01"
    data_final := "2026-04-01"
    servico := financas.SimularJurosComposto {}
    servico.SetValorInicial(0.0)
    servico.SetDatas(data_inicial, data_final)
    servico.SetTaxaAnos(2)
    servico.SetTaxaMeses(servico.GetTaxaAnos() * 12)
    servico.SetDiasDeLiquidesPorAno(12)
    servico.SetTaxaDeJurosDecimal(13.25, financas.PORCENTO_ANUAL)
    servico.SetValorAporte(833.00)
    _, err := financas.FutureValueOfASeriesMonthly(servico.GetValorInicial(), servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), servico.GetValorAporte(), servico.GetTaxaMeses(), true, servico.GetDataInicial())
    if err != nil {
        self.Fail(err.Error())
    }
}

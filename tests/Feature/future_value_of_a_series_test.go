package Feature

import (
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
func (self *FutureValueOfASeriesSuite) TestNewStruct() {
    servico := financas.SimulacaoCDB {}
    servico.SetDatas("2025-01-01", "2026-01-01")
    taxa_selic := servico.GetTaxaSelic()
    servico.SetDiasDeLiquidesPorAno(255)
    servico.SetTaxaDeJurosDecimal(taxa_selic, "porcento anual")
    servico.SetValorAporte(833.00)
    servico.SetDiasDeLiquidesPorAno(12)
    servico.SetTaxaAnosApartirPeriodoDeDatas()
    valorizacao := financas.FutureValuesOfASeriesFormula(servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), 1, servico.GetValorAporte())
    if valorizacao != 10743.284798509316 {
        self.Fail("Valorização deveria ser 10743.284798509316")
    }
}
func (self *FutureValueOfASeriesSuite) TestCompoundInterestFomula() {
    servico := financas.SimulacaoCDB {}
    servico.SetDatas("2025-01-01", "2026-01-01")
    taxa_selic := servico.GetTaxaSelic()
    servico.SetDiasDeLiquidesPorAno(255)
    servico.SetTaxaDeJurosDecimal(taxa_selic, "porcento anual")
    servico.SetValorInicial(9500.00)
    servico.SetDiasDeLiquidesPorAno(12)
    servico.SetTaxaAnosApartirPeriodoDeDatas()
    valorizacao := financas.CompoundInterestFormula(servico.GetValorInicial(), servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), servico.GetTaxaAnos())
    if valorizacao != 10842.00177695996 {
        self.Fail("Valorização deveria ser 10842.00177695996")
    }
}

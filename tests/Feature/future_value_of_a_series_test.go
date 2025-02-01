package Feature

import (
	"fmt"
	"testing"

	"goravel/app/core"
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
    servico := financas.SimularJurosComposto {}
    servico.SetDiasDeLiquidesPorAno(255)
    servico.SetTaxaDeJurosDecimal(13.25, financas.PROCENTO_ANUAL)
    servico.SetValorAporte(833.00)
    servico.SetDiasDeLiquidesPorAno(12)
    valorizacao := financas.FutureValuesOfASeriesFormula(servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), 1, servico.GetValorAporte(), false)
    valor_alvo := 10625.956528507051
    if valorizacao != valor_alvo {
        self.Fail(fmt.Sprintf("[FVS] Valorização deveria ser %f, retornado %f", valor_alvo, valorizacao))
    }
}
func (self *FutureValueOfASeriesSuite) TestCompoundInterestFomula() {
    servico := financas.SimularJurosComposto {}
    servico.SetDiasDeLiquidesPorAno(255)
    servico.SetTaxaDeJurosDecimal(13.25, financas.PROCENTO_ANUAL)
    servico.SetValorInicial(9500.00)
    servico.SetDiasDeLiquidesPorAno(12)
    valorizacao := financas.CompoundInterestFormula(servico.GetValorInicial(), servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), 1)
    valor_alvo := 10838.077509029437
    if valorizacao != valor_alvo {
        self.Fail(fmt.Sprintf("[CI] Valorização deveria ser %f, retornado %f", valor_alvo, valorizacao))
    }
}
func (self *FutureValueOfASeriesSuite) TestCivFvs() {
    servico := financas.SimularJurosComposto {}
    servico.SetTaxaDeJurosDecimal(13.25, financas.PROCENTO_ANUAL)
    servico.SetDiasDeLiquidesPorAno(255)
    servico.SetValorInicial(9500.00)
    servico.SetDiasDeLiquidesPorAno(12)
    servico.SetValorAporte(833.00)
    valorizacao := financas.CifAndFvs(servico.GetValorInicial(), servico.GetTaxaDeJurosDecimal(), servico.GetDiasDeLiquidezPorAno(), 1, servico.GetValorAporte())
    valor_alvo := 21581.362308
    if !core.AlmostEqual(valorizacao, valor_alvo, 0.1) {
        self.Fail(fmt.Sprintf("[CIF FVS] Valorização deveria ser %f, retornado %f", valor_alvo, valorizacao))
    }
}

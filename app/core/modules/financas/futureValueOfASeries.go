package financas

import (
	"github.com/shopspring/decimal"
)

// interestRateDecimal: decimal interest rate per period
// periods: number of days until payment
// months: investment duration in years
// contributionAmount: amount of each contribution
// ContributionOnFirstDay: whether contribution is made at start of period
type FutureValueOfASeries struct {
    InterestRateDecimal float64
    Periods             float64
    ContributionAmount  float64
    ContributionOnFirstDay bool
}
func (self * FutureValueOfASeries) Calculate() float64 {
    decimalContribuitionAmount := decimal.NewFromFloat(self.ContributionAmount)
    decimalPeriods := decimal.NewFromFloat(self.Periods)
    decimalInterestRateDecimal := decimal.NewFromFloat(self.InterestRateDecimal)
    taxPerPeriods := decimalInterestRateDecimal.Round(16).Div(decimalPeriods).Round(16)
    growthFactor := decimal.NewFromInt(1).Add(taxPerPeriods).Round(16).Pow(decimalPeriods).Round(16)
    growthFactorPerTaxPerPeridos := growthFactor.Sub(decimal.NewFromInt(1)).Div(taxPerPeriods)
    result := decimalContribuitionAmount.Mul(growthFactorPerTaxPerPeridos).Round(16)
    if self.ContributionOnFirstDay {
        firstPeriod := decimal.NewFromInt(1).Add(taxPerPeriods)
        futureValue, _ := result.Mul(firstPeriod).Round(16).Float64()
        return futureValue
    }
    futureValue, _ := result.Round(16).Float64()
    return futureValue
}

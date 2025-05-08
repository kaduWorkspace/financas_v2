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
    InterestRateDecimal float64 //Annually
    Periods             float64
    ContributionAmount  float64
    ContributionOnFirstDay bool
}
func (self *FutureValueOfASeries) Calculate() float64 {
    // Convert inputs to decimal for precise calculations
    contrib := decimal.NewFromFloat(self.ContributionAmount)
    periods := decimal.NewFromFloat(self.Periods)
    rate := decimal.NewFromFloat(self.InterestRateDecimal)

    // Calculate periodic interest rate
    periodicRate := rate.Round(16).Div(decimal.NewFromInt(12)).Round(16)

    // Calculate growth factor
    growthFactor := decimal.NewFromInt(1).
        Add(periodicRate).
        Round(16).
        Pow(periods).
        Round(16)

    // Calculate future value factor
    futureValueFactor := growthFactor.
        Sub(decimal.NewFromInt(1)).
        Div(periodicRate)

    // Calculate base result
    result := contrib.Mul(futureValueFactor).Round(16)

    // Adjust for first day contribution if needed
    if self.ContributionOnFirstDay {
        firstPeriodGrowth := decimal.NewFromInt(1).Add(periodicRate)
        result = result.Mul(firstPeriodGrowth).Round(16)
    }

    // Convert final result to float64
    finalValue, _ := result.Float64()
    return finalValue
}
type period struct {
    Accrued float64 `json:"accrued"`
    Period  int `json:"period"`
    Interest float64 `json:"interest"`
}
func (self * FutureValueOfASeries) CalculateWithPeriods(initialValue float64) float64 {
    decimalInitialValue := decimal.NewFromFloat(initialValue)
    decimalContribuitionAmount := decimal.NewFromFloat(self.ContributionAmount)
    decimalInterestRateDecimal := decimal.NewFromFloat(self.InterestRateDecimal).Div(decimal.NewFromInt(12))
    periodsInt := int(self.Periods)
    counter := 0
    accrued := decimal.NewFromInt(0).Add(decimalInitialValue)
    periods := []period{}
    for counter < periodsInt {
        accrued = accrued.Add(decimalContribuitionAmount)
        accruedInterest := decimalInterestRateDecimal.Mul(accrued)
        accrued = accrued.Add(accruedInterest)
        t,_:=accrued.Round(16).Float64()
        i,_ := accruedInterest.Round(16).Float64()
        periods = append(periods, period{
            Accrued: t,
            Period: counter + 1,
            Interest: i,
        })
        counter++
    }
    futureValue, _ := accrued.Round(16).Float64()
    return futureValue
}

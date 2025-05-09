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
    Investment
    interestRateDecimal float64 //Annually
    periods             float64
    contributionAmount  float64
    contributionOnFirstDay bool
}

func (f *FutureValueOfASeries) InterestRateDecimal() float64 {
    return f.interestRateDecimal
}

func (f *FutureValueOfASeries) SetInterestRateDecimal(rate float64) {
    f.interestRateDecimal = rate
}

func (f *FutureValueOfASeries) Periods() float64 {
    return f.periods
}

func (f *FutureValueOfASeries) SetPeriods(periods float64) {
    f.periods = periods
}

func (f *FutureValueOfASeries) ContributionAmount() float64 {
    return f.contributionAmount
}

func (f *FutureValueOfASeries) SetContributionAmount(amount float64) {
    f.contributionAmount = amount
}

func (f *FutureValueOfASeries) ContributionOnFirstDay() bool {
    return f.contributionOnFirstDay
}

func (f *FutureValueOfASeries) SetContributionOnFirstDay(onFirstDay bool) {
    f.contributionOnFirstDay = onFirstDay
}
func (self *FutureValueOfASeries) Calculate() float64 {
    // Convert inputs to decimal for precise calculations
    contrib := decimal.NewFromFloat(self.contributionAmount)
    periods := decimal.NewFromFloat(self.periods)
    rate := decimal.NewFromFloat(self.interestRateDecimal)

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
    if self.contributionOnFirstDay {
        firstPeriodGrowth := decimal.NewFromInt(1).Add(periodicRate)
        result = result.Mul(firstPeriodGrowth).Round(16)
    }

    // Convert final result to float64
    finalValue, _ := result.Float64()
    return finalValue
}
type Period struct {
    Accrued float64 `json:"accrued"`
    Period  int `json:"period"`
    Interest float64 `json:"interest"`
}
func (self * FutureValueOfASeries) CalculateWithPeriods(initialValue float64) (float64, []Period) {
    decimalInitialValue := decimal.NewFromFloat(initialValue)
    decimalContribuitionAmount := decimal.NewFromFloat(self.contributionAmount)
    decimalInterestRateDecimal := decimal.NewFromFloat(self.interestRateDecimal).Div(decimal.NewFromInt(12))
    periodsInt := int(self.periods)
    counter := 0
    accrued := decimal.NewFromInt(0).Add(decimalInitialValue)
    periods := []Period{}
    for counter < periodsInt {
        accrued = accrued.Add(decimalContribuitionAmount)
        accruedInterest := decimalInterestRateDecimal.Mul(accrued)
        accrued = accrued.Add(accruedInterest)
        t,_:=accrued.Round(16).Float64()
        i,_ := accruedInterest.Round(16).Float64()
        periods = append(periods, Period{
            Accrued: t,
            Period: counter + 1,
            Interest: i,
        })
        counter++
    }
    futureValue, _ := accrued.Round(16).Float64()
    return futureValue, periods
}

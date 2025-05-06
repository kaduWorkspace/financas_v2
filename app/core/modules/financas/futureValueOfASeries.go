package financas
// interestRateDecimal: decimal interest rate per period
// netDays: number of days until payment
// years: investment duration in years
// contributionAmount: amount of each contribution
// contributionOnFirstDay: whether contribution is made at start of period
type FutureValueOfASeries struct {
    initialValue float64
    interestRateDecimal float64
    netDays             float64
    years               float64
    contributionAmount  float64
    contributionOnFirstDay bool
}

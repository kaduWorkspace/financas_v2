package financas

import (
	"github.com/shopspring/decimal"
)

type CompoundInterest struct {
    Investment
    InitialValue float64
    Tax  float64
    Months int
}

func (self *CompoundInterest) Calculate() float64 {
    decimalInitialValue := decimal.NewFromFloat(self.InitialValue)
    decimalTax := decimal.NewFromFloat(self.Tax)
    decimalMonths := decimal.NewFromInt(int64(self.Months))
    futureValue := decimalInitialValue.Mul(
        decimal.NewFromInt(1).Add(decimalTax.Div(decimal.NewFromInt(12))).Pow(decimalMonths),
    )
    result, _ := futureValue.Round(16).Float64()
    return result
}


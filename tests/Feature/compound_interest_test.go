package Feature

import (
	"goravel/app/core/modules/financas"
	"math"
	"testing"
)

func TestCompoundInterestCalculate(t *testing.T) {
    tests := []struct {
        name     string
        initial  float64
        tax      float64
        months   int
        expected float64
    }{
        {
            name:     "Small investment with low interest",
            initial:  1000.0,
            tax:      0.01,
            months:   12,
            expected: 1126.8250301319697,
        },
        {
            name:     "Large investment with moderate interest",
            initial:  1000000.0,
            tax:      0.05,
            months:   24,
            expected: 3225099.9437137,
        },
        {
            name:     "Small investment with high interest",
            initial:  500.0,
            tax:      0.2,
            months:   6,
            expected: 1492.992,
        },
        {
            name:     "Large investment with small monthly gain",
            initial:  5000000.0,
            tax:      0.001,
            months:   36,
            expected: 5183185.996419741,
        },
        {
            name:     "Zero months should return initial value",
            initial:  1000.0,
            tax:      0.1,
            months:   0,
            expected: 1000.0,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ci := &financas.CompoundInterest{
                InitialValue: tt.initial,
                Tax:          tt.tax,
                Months:       tt.months,
            }
            got := ci.Calculate()
            tolerance := 1e-8 // 0.000001% tolerance
            diff := math.Abs(got - tt.expected)
            if diff > tolerance && diff/math.Max(math.Abs(got), math.Abs(tt.expected)) > tolerance {
                t.Errorf("Calculate() = %v, want %v (diff: %v)", got, tt.expected, diff)
            }
        })
    }
}

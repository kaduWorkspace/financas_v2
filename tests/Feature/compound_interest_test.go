package Feature

import (
	"fmt"
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
            expected: 1010.045,
        },
        {
            name:     "Large investment with moderate interest",
            initial:  10000.0,
            tax:      0.05,
            months:   24,
            expected: 11049.41,
        },
        {
            name:     "Small investment with high interest",
            initial:  500.0,
            tax:      0.2,
            months:   6,
            expected: 552.13,
        },
        {
            name:     "Large investment with small monthly gain",
            initial:  50000.0,
            tax:      0.001,
            months:   36,
            expected: 50150.21,
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
            tolerance := 0.01 // 0.000001% tolerance
            diff := math.Abs(got - tt.expected)
            if diff > tolerance && diff/math.Max(math.Abs(got), math.Abs(tt.expected)) > tolerance {
                t.Errorf("Calculate() = %v, want %v (diff: %v)", got, tt.expected, diff)
            }
        })
    }
}

func TestFutureValueOfASeries(t *testing.T) {
	tests := []struct {
		name                  string
		interestRateDecimal   float64
		periods               float64
		contributionAmount    float64
		contributionOnFirstDay bool
		want                  float64
	}{
		{
			name:                  "monthly contributions end period",
			interestRateDecimal:   0.12, // 12% annual
			periods:               12,   // 1 year monthly
			contributionAmount:    100,
			contributionOnFirstDay: false,
			want:                  1268.25,
		},
		{
			name:                  "monthly contributions start period",
			interestRateDecimal:   0.12,
			periods:               12,
			contributionAmount:    100,
			contributionOnFirstDay: true,
			want:                  1280.93,
		},
		{
			name:                  "quarterly contributions low rate",
			interestRateDecimal:   0.01, // 1% annual
			periods:               4,    // 4 months
			contributionAmount:    500,
			contributionOnFirstDay: false,
			want:                  2002.50,
		},
		{
			name:                  "weekly contributions high rate",
			interestRateDecimal:   0.24, // 24% annual
			periods:               52,   // 52 months 4 anos e uns 4 meses
			contributionAmount:    10,
			contributionOnFirstDay: true,
			want:                  918.16,
		},
		{
			name:                  "single contribution edge case",
			interestRateDecimal:   0.05,
			periods:               1,
			contributionAmount:    1000,
			contributionOnFirstDay: false,
			want:                  1000.00,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := financas.FutureValueOfASeries{
				InterestRateDecimal:   tt.interestRateDecimal,
				Periods:               tt.periods,
				ContributionAmount:    tt.contributionAmount,
				ContributionOnFirstDay: tt.contributionOnFirstDay,
			}
			got := fv.Calculate()
			if !almostEqual(got, tt.want, 0.01) {
				t.Errorf("calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestCompareFormulaWithLoop(t *testing.T) {
    data := struct {
        name                  string
        interestRateDecimal   float64
        periods               float64
        contributionAmount    float64
        contributionOnFirstDay bool
    }{
        name:                  "single contribution edge case",
        interestRateDecimal:   0.1425,
        periods:               36,
        contributionAmount:    1000,
        contributionOnFirstDay: true,
    }
    fv := financas.FutureValueOfASeries{
        InterestRateDecimal:   data.interestRateDecimal,
        Periods:               data.periods,
        ContributionAmount:    data.contributionAmount,
        ContributionOnFirstDay: data.contributionOnFirstDay,
    }
    one := fv.Calculate()
    two := fv.CalculateWithPeriods(0)
    if !almostEqual(one, two, 0.01) {
        t.Errorf("one = %v, want %v", one, two)
    }
    fmt.Println(one, two)
}
func TestCompareFormulaWithLoop2(t *testing.T) {
    data := struct {
        name                  string
        interestRateDecimal   float64
        periods               float64
        contributionAmount    float64
        contributionOnFirstDay bool
        want                  float64
    }{
        name:                  "single contribution edge case",
        interestRateDecimal:   0.1425,
        periods:               10,
        contributionAmount:    1000,
        contributionOnFirstDay: true,
        want:                  1000.00,
    }
    fv := financas.FutureValueOfASeries{
        InterestRateDecimal:   data.interestRateDecimal,
        Periods:               data.periods,
        ContributionAmount:    data.contributionAmount,
        ContributionOnFirstDay: data.contributionOnFirstDay,
    }
    cp := financas.CompoundInterest{
        InitialValue: 100,
        Tax: 0.1425,
        Months: 10,
    }
    one := fv.Calculate()
    one_cp := cp.Calculate()
    two := fv.CalculateWithPeriods(100)
    if !almostEqual(one + one_cp, two, 0.0001) {
        t.Errorf("one = %v, want %v", one + one_cp, two)
    }
    fmt.Println(one, one_cp, two)
}

// Helper function to compare floating point numbers with tolerance
func almostEqual(a, b, tolerance float64) bool {
	return (a-b) < tolerance && (b-a) < tolerance
}

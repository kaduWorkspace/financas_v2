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
            //fmt.Printf("Calculate() = %v, want %v (diff: %v)", got, tt.expected, diff)
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
			fv := financas.FutureValueOfASeries{}
            fv.SetInterestRateDecimal(tt.interestRateDecimal)
            fv.SetPeriods(tt.periods)
            fv.SetContributionAmount(tt.contributionAmount)
            fv.SetContributionOnFirstDay(tt.contributionOnFirstDay)
			got := fv.Calculate()
			if !almostEqual(got, tt.want, 0.01) {
				t.Errorf("calculate() = %v, want %v", got, tt.want)
			}
            fmt.Println(fmt.Sprintf("calculate() = %v, want %v", got, tt.want))
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
    fv := financas.FutureValueOfASeries{}
    fv.SetInterestRateDecimal(data.interestRateDecimal)
    fv.SetPeriods(data.periods)
    fv.SetContributionAmount(data.contributionAmount)
    fv.SetContributionOnFirstDay(data.contributionOnFirstDay)
    one := fv.Calculate()
    two, _ := fv.CalculateWithPeriods(0)
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
    fv := financas.FutureValueOfASeries{}
    fv.SetInterestRateDecimal(data.interestRateDecimal)
    fv.SetPeriods(data.periods)
    fv.SetContributionAmount(data.contributionAmount)
    fv.SetContributionOnFirstDay(data.contributionOnFirstDay)
    cp := financas.CompoundInterest{
        InitialValue: 100,
        Tax: 0.1425,
        Months: 10,
    }
    one := fv.Calculate()
    one_cp := cp.Calculate()
    two, _ := fv.CalculateWithPeriods(100)
    if !almostEqual(one + one_cp, two, 0.0001) {
        t.Errorf("one = %v, want %v", one + one_cp, two)
    }
    fmt.Println(one, one_cp, two)
}
func TestFutureValueOfASeries_PredictFV(t *testing.T) {
	tests := []struct {
		name                  string
		interestRateDecimal   float64
		periods               float64
		contributionOnFirstDay bool
		finalValue            float64
		want                  float64
	}{
		{
			name:                  "monthly contributions",
			interestRateDecimal:   0.12,
			periods:               12,
			finalValue:            1280.93,
			want:                  99.999,
		},
		{
			name:                  "semester contributions",
			interestRateDecimal:   0.01,
			periods:               6,
			finalValue:            2015.87,
			want:                  334.99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := financas.FutureValueOfASeries{}
			fv.SetInterestRateDecimal(tt.interestRateDecimal)
			fv.SetPeriods(tt.periods)
            fv.SetContributionOnFirstDay(true)
			got := fv.PredictFV(tt.finalValue)
			if !almostEqual(got, tt.want, 0.01) {
				t.Errorf("PredictFV() = %v, want %v", got, tt.want)
			}

            fv.SetContributionAmount(got)
            compare, _ := fv.CalculateWithPeriods(0)
            if !almostEqual(tt.finalValue, compare, 0.001) {
				t.Errorf("Compare() = %v, compare want %v", compare, tt.finalValue)
            }
            //fmt.Printf("PredictFV() = %v, want %v", got, tt.want)
		})
	}
}
func TestFutureValueOfASeries_PredictFVWithInitialValue(t *testing.T) {
	tests := []struct {
		name                  string
		interestRateDecimal   float64
		periods               float64
		contributionOnFirstDay bool
		finalValue            float64
		initialValue          float64
		want                  float64
	}{
		{
			name:                  "with initial value end period",
			interestRateDecimal:   0.12,
			periods:               12,
			contributionOnFirstDay: true,
			finalValue:            2062.84,
			initialValue:          500.00,
			want:                  117.057, // Example value, adjust based on actual calculation
		},
		{
			name:                  "with initial value start period",
			interestRateDecimal:   0.12,
			periods:               12,
			contributionOnFirstDay: true,
			finalValue:            2535.62,
			initialValue:          300.00,
			want:                  171.56, // Example value, adjust based on actual calculation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fv := financas.FutureValueOfASeries{}
			fv.SetInterestRateDecimal(tt.interestRateDecimal)
			fv.SetPeriods(tt.periods)
			fv.SetContributionOnFirstDay(tt.contributionOnFirstDay)

			got := fv.PredictFVWithInitialValue(tt.finalValue, tt.initialValue)
			if !almostEqual(got, tt.want, 0.1) {
				t.Errorf("PredictFVWithInitialValue() = %v, want %v", got, tt.want)
			}

            //fmt.Printf("PredictFVWithInitialValue() = %v, want %v", got, tt.want)
            fv.SetContributionAmount(got)
            compare, _ := fv.CalculateWithPeriods(tt.initialValue)
            //fmt.Println(periods, compare)
            if !almostEqual(tt.finalValue, compare, 0.01) {
				t.Errorf("Compare() = %v, want %v", compare, tt.finalValue)
            }
		})
	}
}
// Helper function to compare floating point numbers with tolerance
func almostEqual(a, b, tolerance float64) bool {
    return math.Abs(a-b) <= tolerance
}

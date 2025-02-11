package calculation

import (
	"testing"
)

func TestGetHeikinDailyClose(t *testing.T) {
	cases := []struct {
		input    PriceCal
		expected int
	}{{
		input: PriceCal{
			Open:  162.96,
			Close: 160.84,
			High:  163.4,
			Low:   158.58,
		},
		expected: 16144,
	}}

	for _, c := range cases {
		price := GetHeikinDailyClose(&c.input)

		if int(price*100) != c.expected {
			t.Fatalf("Expected price: %f, actual price: %f", float64(c.expected)/100, price)
		}
	}
}

func TestGetHeikinDailyLow(t *testing.T) {
	cases := []struct {
		input    PriceCal
		expected int
	}{{
		input: PriceCal{
			Open:  162.96,
			Close: 160.84,
			High:  163.4,
			Low:   158.58,
		},
		expected: 15858,
	}}

	for _, c := range cases {
		price := GetHeikinDailyLow(&c.input)

		if int(price*100) != c.expected {
			t.Fatalf("Expected price: %f, actual price: %f", float64(c.expected)/100, price)
		}
	}
}

func TestGetHeikinDailyHigh(t *testing.T) {
	cases := []struct {
		input    PriceCal
		expected int
	}{{
		input: PriceCal{
			Open:  162.96,
			Close: 160.84,
			High:  163.4,
			Low:   158.58,
		},
		expected: 16340,
	}}

	for _, c := range cases {
		price := GetHeikinDailyHigh(&c.input)

		if int(price*100) != c.expected {
			t.Fatalf("Expected price: %f, actual price: %f", float64(c.expected)/100, price)
		}
	}
}

func TestGetHeikinDailyOpen(t *testing.T) {
	cases := []struct {
		prev     PriceCal
		expected int
	}{{
		prev: PriceCal{
			Open:  163.69,
			Close: 165.15,
			High:  170.74,
			Low:   160.87,
		},
		expected: 16442,
	}}

	for _, c := range cases {
		price := GetHeikinDailyOpen(&c.prev)

		if int(price*100) != c.expected {
			t.Fatalf("Expected price: %f, actual price: %f", float64(c.expected)/100, price)
		}
	}
}

func TestIsDojiStar(t *testing.T) {
	cases := []struct {
		input    PriceCal
		prev     PriceCal
		expected bool
	}{
		{ // ARM 2024-08-06
			input: PriceCal{
				Open:  115.53,
				High:  117.97,
				Low:   109.50,
				Close: 113.39,
			},
			prev: PriceCal{
				Open:  123.73,
				Close: 104.78,
			},
			expected: true,
		},
		{ // ARM 2024-12-09
			input: PriceCal{
				Open:  140.15,
				High:  143.20,
				Low:   136.26,
				Close: 139.64,
			},
			prev: PriceCal{
				Open:  140.01,
				Close: 139.52,
			},
			expected: true,
		},
		{ // ARM 2024-10-16
			input: PriceCal{
				Open:  154.00,
				High:  155.20,
				Low:   151.29,
				Close: 152.50,
			},
			prev: PriceCal{
				Open:  153.08,
				Close: 154.57,
			},
			expected: false,
		},
		{ // ARM 2024-10-15
			input: PriceCal{
				Open:  160.00,
				High:  160.62,
				Low:   147.00,
				Close: 150.67,
			},
			prev: PriceCal{
				Open:  148.11,
				Close: 158.05,
			},
			expected: false,
		},
		{ // ARM 2024-10-14
			input: PriceCal{
				Open:  153.10,
				High:  164.16,
				Low:   153.10,
				Close: 161.82,
			},
			prev: PriceCal{
				Open:  145.94,
				Close: 150.29,
			},
			expected: false,
		},
	}

	for i, c := range cases {
		isDoji := GetIsDojiStar(&c.input, &c.prev)

		if isDoji != c.expected {
			t.Fatalf("Test Case %d of TestIsDojiStar failed, expected: %t, actual: %t", i, c.expected, isDoji)
		}
	}
}

func TestIsSpinningTop(t *testing.T) {
	cases := []struct {
		input    PriceCal
		prev     PriceCal
		expected bool
	}{
		{ // ARM 2024-12-09
			input: PriceCal{
				Open:  140.15,
				High:  143.20,
				Low:   136.26,
				Close: 139.64,
			},
			prev: PriceCal{
				Open:  140.01,
				Close: 139.52,
			},
			expected: false,
		},
		{ // ARM 2024-10-16
			input: PriceCal{
				Open:  154.00,
				High:  155.20,
				Low:   151.29,
				Close: 152.50,
			},
			prev: PriceCal{
				Open:  153.08,
				Close: 154.57,
			},
			expected: true,
		},
		{ // ARM 2024-10-15
			input: PriceCal{
				Open:  160.00,
				High:  160.62,
				Low:   147.00,
				Close: 150.67,
			},
			prev: PriceCal{
				Open:  148.11,
				Close: 158.05,
			},
			expected: true,
		},
		{ // ARM 2024-10-14
			input: PriceCal{
				Open:  153.10,
				High:  164.16,
				Low:   153.10,
				Close: 161.82,
			},
			prev: PriceCal{
				Open:  145.94,
				Close: 150.29,
			},
			expected: false,
		},

		{ // ARM 2024-10-18
			input: PriceCal{
				Open:  155.57,
				High:  155.74,
				Low:   151.96,
				Close: 153.03,
			},
			prev: PriceCal{
				Open:  153.54,
				Close: 156.30,
			},
			expected: false,
		},

		{ // ARM 2024-10-22
			input: PriceCal{
				Open:  150.46,
				High:  152.94,
				Low:   149.83,
				Close: 152.58,
			},
			prev: PriceCal{
				Open:  154.50,
				Close: 152.01,
			},
			expected: false,
		},
		{ // ARM 2024-10-25
			input: PriceCal{
				Open:  142.00,
				High:  145.56,
				Low:   141.50,
				Close: 143.75,
			},
			prev: PriceCal{
				Open:  148.60,
				Close: 141.58,
			},
			expected: false,
		},
	}

	for i, c := range cases {
		isSpinningTop := GetIsSpinningTop(&c.input, &c.prev)

		if isSpinningTop != c.expected {
			t.Fatalf("Test Case %d of TestIsSpinningTop failed, expected: %t, actual: %t", i, c.expected, isSpinningTop)
		}
	}
}

func TestGetIsUptrend(t *testing.T) {
	cases := []struct {
		input    PriceCal
		prev     PriceCal
		expected bool
	}{
		{ // ARM 2024-12-09
			input: PriceCal{
				Open:  140.15,
				High:  143.20,
				Low:   136.26,
				Close: 139.64,
			},
			prev: PriceCal{
				Open:  140.01,
				Close: 139.52,
			},
			expected: true,
		},
		{ // ARM 2024-10-16
			input: PriceCal{
				Open:  154.00,
				High:  155.20,
				Low:   151.29,
				Close: 152.50,
			},
			prev: PriceCal{
				Open:  153.08,
				Close: 154.57,
			},
			expected: false,
		},
		{ // ARM 2024-10-15
			input: PriceCal{
				Open:  160.00,
				High:  160.62,
				Low:   147.00,
				Close: 150.67,
			},
			prev: PriceCal{
				Open:  148.11,
				Close: 158.05,
			},
			expected: true,
		},
		{ // ARM 2024-10-14
			input: PriceCal{
				Open:  153.10,
				High:  164.16,
				Low:   153.10,
				Close: 161.82,
			},
			prev: PriceCal{
				Open:  145.94,
				Close: 150.29,
			},
			expected: true,
		},

		{ // ARM 2024-10-18
			input: PriceCal{
				Open:  155.57,
				High:  155.74,
				Low:   151.96,
				Close: 153.03,
			},
			prev: PriceCal{
				Open:  153.54,
				Close: 156.30,
			},
			expected: false,
		},

		{ // ARM 2024-10-22
			input: PriceCal{
				Open:  150.46,
				High:  152.94,
				Low:   149.83,
				Close: 152.58,
			},
			prev: PriceCal{
				Open:  154.50,
				Close: 152.01,
			},
			expected: false,
		},
		{ // ARM 2024-10-25
			input: PriceCal{
				Open:  142.00,
				High:  145.56,
				Low:   141.50,
				Close: 143.75,
			},
			prev: PriceCal{
				Open:  148.60,
				Close: 141.58,
			},
			expected: false,
		},
	}

	for i, c := range cases {
		isUptrend := GetIsUptrend(&c.input, &c.prev)

		if isUptrend != c.expected {
			t.Fatalf("Test Case %d of TestGetIsUptrend failed, expected: %t, actual: %t", i, c.expected, isUptrend)
		}
	}
}

package divide

import (
	"testing"

	"github.com/the-anna-project/context"
)

func Test_Service_Action(t *testing.T) {
	testCases := []struct {
		A        float64
		B        float64
		Expected float64
	}{
		{
			A:        3.5,
			B:        12.5,
			Expected: 0.28,
		},
		{
			A:        35.5,
			B:        14.5,
			Expected: 2.4482758620689653,
		},
		{
			A:        -3.5,
			B:        7.5,
			Expected: -0.4666666666666667,
		},
		{
			A:        12.5,
			B:        4.5,
			Expected: 2.7777777777777777,
		},
		{
			A:        36.5,
			B:        6.5,
			Expected: 5.615384615384615,
		},
	}

	newService, err := NewService(DefaultServiceConfig())
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	for i, testCase := range testCases {
		action := newService.Action().(func(ctx context.Context, a, b float64) float64)
		f := action(nil, testCase.A, testCase.B)
		if f != testCase.Expected {
			t.Fatal("case", i+1, "expected", testCase.Expected, "got", f)
		}
	}
}

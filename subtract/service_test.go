package subtract

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
			Expected: -9,
		},
		{
			A:        35.5,
			B:        14.5,
			Expected: 21,
		},
		{
			A:        -3.5,
			B:        -7.5,
			Expected: 4,
		},
		{
			A:        12.5,
			B:        4.5,
			Expected: 8,
		},
		{
			A:        36.5,
			B:        6.5,
			Expected: 30,
		},
		{
			A:        11.11,
			B:        10.10,
			Expected: 1.0099999999999998,
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
